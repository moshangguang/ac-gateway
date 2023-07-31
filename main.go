package main

import (
	"ac-gateway/config"
	"ac-gateway/help/log"
	"ac-gateway/initialize"
	"ac-gateway/pkg/di"
	"ac-gateway/pkg/server"
	"go.uber.org/zap"
)

func main() {
	initialize.Init()
	if err := StartServer(); err != nil {
		log.Fatal("服务启动失败", zap.Error(err))
	}
}
func StartServer() error {
	return di.Container.Invoke(func(server *server.Server, config config.Config) error {
		return server.Start(config)
	})
}
