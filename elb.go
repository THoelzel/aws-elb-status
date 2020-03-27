package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"go.uber.org/zap"
	"strings"
)

var (
	sess *session.Session
	instanceNames map[string]string
)

type status interface {
	refresh()
}

type Status struct {
	Name     string
	Clusters []Cluster
}

type MockStatus struct {
	Name     string
	Clusters []Cluster
}

// Represents LoadBalancers
type Cluster struct {
	Name             string
	LoadBalancerArn  string
	ServerCategories []ServerCategory
	Status           int
}

// Represents TargetGroups
type ServerCategory struct {
	Name           string
	Servers        []Server
	TargetGroupArn string
	Status         int
}

type Server struct {
	Name   string
	Status int
}

func NewStatus(name string) Status {

	sess = session.Must(session.NewSession())
	instanceNames = nameAllInstances()

	return Status{
		Name:     name,
		Clusters: findClusters(),
	}
}

func (c *Status) refresh() {

}

func (m *MockStatus) refresh() {

}

func findClusters() []Cluster {

	var clusters []Cluster

	in := elbv2.DescribeLoadBalancersInput{}
	svc := elbv2.New(sess)

	out, err := svc.DescribeLoadBalancers(&in)

	if err != nil {
		logger.Error("unable to fetch load balancers",
			zap.Error(err))
		return clusters
	}

	for _, l := range out.LoadBalancers {

		c := Cluster{
			Name:            *l.LoadBalancerName,
			LoadBalancerArn: *l.LoadBalancerArn,
		}
		clusters = append(clusters, c)
	}

	logger.Debug("located load balancers",
		zap.Int("count", len(clusters)))

	return findAllTargetGroups(clusters)
}

func findAllTargetGroups(clusters []Cluster) []Cluster {
	in := elbv2.DescribeTargetGroupsInput{}

	// TODO
	svc := elbv2.New(sess)

	out, err := svc.DescribeTargetGroups(&in)

	if err != nil {
		logger.Error("unable to fetch target groups",
			zap.Error(err))
	}

	labeledClusters := make(map[string]Cluster)

	for _, c := range clusters {
		labeledClusters[c.LoadBalancerArn] = c
		// TODO
	}

	if out != nil {
		logger.Debug("located target groups",
			zap.Int("count", len(out.TargetGroups)))

		// only supports single LB pointing to TG
		for _, t := range out.TargetGroups {
			if len(t.LoadBalancerArns) > 0 {
				if v, ok := labeledClusters[*t.LoadBalancerArns[0]]; ok {
					s := ServerCategory{
						Name:           *t.TargetGroupName,
						TargetGroupArn: *t.TargetGroupArn,
					}
					s.instanceHealth()
					v.ServerCategories = append(v.ServerCategories, s)
					labeledClusters[*t.LoadBalancerArns[0]] = v
				}
			}
		}
	}

	var populatedClusters []Cluster

	for _, l := range labeledClusters {

		populatedClusters = append(populatedClusters, l)
	}

	return populatedClusters
}

func (s *ServerCategory) instanceHealth() {

	in := elbv2.DescribeTargetHealthInput{
		TargetGroupArn: &s.TargetGroupArn,
	}

	svc := elbv2.New(sess)
	out, err := svc.DescribeTargetHealth(&in)

	if err != nil {
		logger.Error("unable to fetch instances",
			zap.String("target_group_name", s.Name),
			zap.String("target_group_arn", s.TargetGroupArn),
			zap.Error(err))
	} else {
		var servers []Server

		for _, t := range out.TargetHealthDescriptions {
			s := Server{
				Name:   *t.Target.Id,
				Status: convertState(*t.TargetHealth.State),
			}
			servers = append(servers, s)
		}
		s.Servers = servers
	}
}

func nameAllInstances() map[string]string {

	names := make(map[string]string)
	i := ec2.DescribeInstancesInput{}
	svc := ec2.New(sess)
	out, err := svc.DescribeInstances(&i)

	if err != nil {
		logger.Error("unable to fetch instance names",
			zap.Error(err))
	} else {
		for _, i := range out.Reservations[0].Instances {
			for _, t := range i.Tags {
				if strings.EqualFold(*t.Key, "name") {
					names[*i.InstanceId] = *t.Value
					break
				}
			}
		}
	}

	return names
}

func convertState(s string) int {
	// TODO
	return 0
}

func NewMockStatus(name string) Status {
	return Status{
		Name: name,
		Clusters: []Cluster{
			{
				Name: "cluster alpha",
				ServerCategories: []ServerCategory{
					{
						Name: "category A",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
					{
						Name: "category B",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
					{
						Name: "category C",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 2,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         2,
					},
					{
						Name: "category D",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
							{
								Name:   "server 2",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
				},
				Status: 2,
			},
			{
				Name: "cluster beta",
				ServerCategories: []ServerCategory{
					{
						Name: "category A",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
					{
						Name: "category B",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
				},
				Status: 0,
			},
			{
				Name: "cluster gamma",
				ServerCategories: []ServerCategory{
					{
						Name: "category A",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
					{
						Name: "category B",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
				},
				Status: 0,
			},
			{
				Name: "cluster delta",
				ServerCategories: []ServerCategory{
					{
						Name: "category A",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
					{
						Name: "category B",
						Servers: []Server{
							{
								Name:   "server 0",
								Status: 0,
							},
							{
								Name:   "server 1",
								Status: 0,
							},
						},
						TargetGroupArn: "",
						Status:         0,
					},
				},
				Status: 0,
			},
		},
	}
}
