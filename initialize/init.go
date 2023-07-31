package initialize

import (
	"ac-gateway/config"
	"ac-gateway/pkg/admin/controller"
	"ac-gateway/pkg/di"
	"ac-gateway/pkg/etcd"
	"ac-gateway/pkg/event"
	"ac-gateway/pkg/models/dml"
	"ac-gateway/pkg/server"
)

func Init() {
	InitComponent()
	InitDML()
	InitAdminController()
	InitSubscriber()
}

func InitComponent() {
	di.MustProvide(config.NewConfig)
	di.MustProvide(etcd.NewClient)
	di.MustProvide(server.NewServer)
	di.MustProvide(event.NewBus)
}

func InitDML() {
	di.MustProvide(dml.NewRouterEtcdModel)
}

func InitAdminController() {
	di.MustProvide(controller.NewRouterCtrl)
}
func InitSubscriber() {
	di.MustInvoke(server.NewRouterSubscriber)
}
