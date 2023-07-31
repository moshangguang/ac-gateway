package etcd

import (
	"ac-gateway/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func NewClient(config config.Config) *clientv3.Client {
	conf := clientv3.Config{
		Endpoints:   config.EtcdConfig.Endpoints,
		DialTimeout: time.Duration(config.EtcdConfig.DialTimeout) * time.Millisecond,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		panic(err)
	}
	return client
}
