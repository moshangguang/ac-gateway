package dml

import (
	"ac-gateway/constant"
	"ac-gateway/help/log"
	"ac-gateway/help/panicutils"
	"ac-gateway/pkg/event"
	"ac-gateway/pkg/gosafe"
	"ac-gateway/pkg/models/ddl"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sort"
	"time"
)

type RouterEtcdModel struct {
	bus    event.Bus
	client *clientv3.Client
}

func (m *RouterEtcdModel) Add(ctx context.Context, route ddl.Router) (ddl.Router, error) {
	//TODO implement me
	panic("implement me")
}

func (m *RouterEtcdModel) GenRouteId(id int64) string {
	return fmt.Sprintf("%s::%d", constant.AcEtcdRouteInfo, id)
}
func (m *RouterEtcdModel) GetAll(ctx context.Context) (list []ddl.Router, err error) {
	response, err := m.client.Get(ctx, constant.AcEtcdRouteInfo, clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, v := range response.Kvs {
		d := ddl.Router{}
		if err := json.Unmarshal(v.Value, &d); err != nil {
			//todo:ding
			continue
		}
		list = append(list, d)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Priority < list[j].Priority {
			return true
		}
		if list[i].Priority > list[j].Priority {
			return false
		}
		return list[i].Id < list[j].Id
	})
	return
}

func (m *RouterEtcdModel) GetById(ctx context.Context, id int64) (router ddl.Router, exists bool, err error) {
	key := m.GenRouteId(id)
	response, err := m.client.Get(ctx, key)
	if err != nil {
		return
	}
	if len(response.Kvs) == 0 {
		return
	}
	kv := response.Kvs[0]
	if err = json.Unmarshal(kv.Value, &router); err != nil {
		return
	}
	exists = true
	return
}
func (m *RouterEtcdModel) PublishUpdated(ctx context.Context) error {
	_, err := m.client.Put(ctx, constant.AcEtcdRouteUpdated, cast.ToString(time.Now().Unix()))
	return err
}
func (m *RouterEtcdModel) Delete(ctx context.Context, id int64) error {
	key := m.GenRouteId(id)
	response, err := m.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	if response.Deleted == 0 {
		return nil
	}
	return m.PublishUpdated(ctx)

}

func (m *RouterEtcdModel) Save(ctx context.Context, router ddl.Router) (ddl.Router, error) {
	key := m.GenRouteId(router.Id)
	routerBytes, err := json.Marshal(router)
	if err != nil {
		return router, err
	}
	if _, err = m.client.Put(ctx, key, string(routerBytes)); err != nil {
		return router, err
	}
	if err = m.PublishUpdated(ctx); err != nil {
		return router, err
	}
	return router, nil
}
func (m *RouterEtcdModel) ListenUpdated(ctx context.Context) error {
	defer func() {
		if p := recover(); p != nil {
			panicStack := panicutils.GetPanicStack()
			log.Error("监听路由变更事件出现异常", zap.Any("p", p), zap.String("panic_stack", panicStack))
		}
	}()
	watch := m.client.Watch(ctx, constant.AcEtcdRouteUpdated)
	<-watch
	log.Logger.Info("监听到路由变化，更新路由")
	return m.bus.Publish(ctx, event.Event{
		Topic: constant.TopicRouterUpdated,
		Data:  struct{}{},
	})
}
func NewRouterEtcdModel(
	bus event.Bus,
	client *clientv3.Client,
) RouterModel {
	m := &RouterEtcdModel{
		bus:    bus,
		client: client,
	}
	gosafe.GoSafe(func() {
		for {
			if err := m.ListenUpdated(context.Background()); err != nil {
				log.Logger.Error("监控路由变化出错", zap.Error(err))
			}
		}
	})
	return m
}
