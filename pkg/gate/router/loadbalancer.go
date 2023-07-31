package router

import (
	"ac-gateway/constant"
	"ac-gateway/help/contextutils"
	"ac-gateway/help/strutils"
	"context"
	"sync/atomic"
)

// 负载均衡接口
type LoadBalancer interface {
	GetNode(ctx context.Context) (string, error)
}

type IPHashLoadBalancer struct {
	nodes []string
}

func (lb *IPHashLoadBalancer) GetNode(ctx context.Context) (string, error) {
	ip, _ := contextutils.GetStringValue(ctx, constant.CtxRequestIP)
	if strutils.IsNotEmpty(ip) {
		return "", constant.ErrCtxNotFoundIP
	}
	hash := strutils.Hash(ip)
	index := int(hash % uint32(len(lb.nodes)))
	return lb.nodes[index], nil
}

func NewIPHashLoadBalancer(nodes []string) LoadBalancer {
	return &IPHashLoadBalancer{
		nodes: nodes,
	}
}

// 轮询均衡器实现
type RoundRobinLoadBalancer struct {
	nodes []string
	index *int64
}

func (lb *RoundRobinLoadBalancer) GetNode(_ context.Context) (string, error) {
	if len(lb.nodes) == 1 {
		return lb.nodes[0], nil
	}
	newIndex := atomic.AddInt64(lb.index, 1)
	return lb.nodes[int(newIndex%int64(len(lb.nodes)))], nil
}

func (lb *RoundRobinLoadBalancer) ChooseByAddresses(address []string) (string, bool) {
	if len(address) == 0 {
		return "", false
	}
	if len(address) == 1 {
		return address[0], true
	}
	newIndex := atomic.AddInt64(lb.index, 1)
	return address[int(newIndex%int64(len(address)))], true
}

func NewRoundRobinLoadBalancer(nodes []string) LoadBalancer {
	return &RoundRobinLoadBalancer{
		nodes: nodes,
		index: new(int64),
	}
}
