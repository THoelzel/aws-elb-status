package main

import (
	"github.com/aws/aws-sdk-go/service/elbv2"
	"go.uber.org/zap"
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

	in := elbv2.DescribeLoadBalancersInput{}

	// TODO
	svc := elbv2.New(nil)

	out, err := svc.DescribeLoadBalancers(&in)

	if err != nil {

		return nil
	}

	var clusters []Cluster

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
	svc := elbv2.New(nil)

	out, err := svc.DescribeTargetGroups(&in)

	if err != nil {

	}

	labeledClusters := make(map[string]Cluster)

	for _, c := range clusters {
		labeledClusters[c.LoadBalancerArn] = c
	}

	if out != nil {
		logger.Debug("located target groups",
			zap.Int("count", len(out.TargetGroups)))

		// only supports single LB pointing to TG
		for _, t := range out.TargetGroups {
			if v, ok := labeledClusters[*t.LoadBalancerArns[0]]; ok {
				s := ServerCategory{
					Name:           *t.TargetGroupName,
					TargetGroupArn: *t.TargetGroupArn,
				}
				s.instanceHealth()
				v.ServerCategories = append(v.ServerCategories, s)
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

	svc := elbv2.New(nil)
	out, err := svc.DescribeTargetHealth(&in)

	if err != nil {

	}

	var servers []Server

	for _, t := range out.TargetHealthDescriptions {
		s := Server{
			Name:   *t.Target.Id,
			Status: convertState(*t.TargetHealth.State),
		}
		servers = append(servers, s)
	}
}

func nameAllInstances() {
	// TODO
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
