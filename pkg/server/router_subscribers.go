package server

import (
	"ac-gateway/constant"
	"ac-gateway/pkg/event"
	"context"
)

type RouterSubscriber struct {
	server *Server
}

func (sub *RouterSubscriber) Topic() event.Topic {
	return constant.TopicRouterUpdated
}

func (sub *RouterSubscriber) Handle(ctx context.Context, _ event.Event) error {
	return sub.server.LoadRouter(ctx)
}

func NewRouterSubscriber(bus event.Bus, server *Server) error {
	return bus.Subscribe(&RouterSubscriber{
		server: server,
	})
}
