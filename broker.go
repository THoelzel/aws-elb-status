package main

import (
	"go.uber.org/zap"
)

var b *Broker

type Broker struct {
	publishCh chan interface{}
	subCh     chan chan interface{}
	unsubCh   chan chan interface{}
}

func NewBroker() *Broker {
	return &Broker{
		publishCh: make(chan interface{}, 1),
		subCh:     make(chan chan interface{}, 1),
		unsubCh:   make(chan chan interface{}, 1),
	}
}

func (b *Broker) Start() {
	subs := map[chan interface{}]struct{}{}
	for {
		select {
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}

			logger.Debug("New subscriber",
				zap.Int("subscriber_count", len(subs)),
			)
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
			logger.Debug("Removed subscriber",
				zap.Int("subscriber_count", len(subs)),
			)
		case msg := <-b.publishCh:
			logger.Debug("Publishing message",
				zap.Int("subscriber_count", len(subs)),
			)
			for msgCh := range subs {
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Subscribe() chan interface{} {
	msgCh := make(chan interface{}, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh chan interface{}) {
	b.unsubCh <- msgCh
}

func (b *Broker) Publish(msg interface{}) {
	b.publishCh <- msg
}
