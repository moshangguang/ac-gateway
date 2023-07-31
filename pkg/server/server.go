package server

import (
	"ac-gateway/config"
	"ac-gateway/constant"
	"ac-gateway/help/fasthttputils"
	"ac-gateway/help/log"
	"ac-gateway/pkg/admin/middleware"
	"ac-gateway/pkg/admin/route"
	"ac-gateway/pkg/gate/router"
	"ac-gateway/pkg/models/dml"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	// fasthttp对象
	adminEngine   *gin.Engine
	routerModel   dml.RouterModel
	gateServer    *fasthttp.Server
	gateOnce      *sync.Once
	routerMatcher *atomic.Value
	etcdClient    *clientv3.Client
	fastClient    *fasthttp.Client
}

func (srv *Server) Health(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	_, _ = ctx.WriteString("You Build It, You Run It.")
}
func (srv *Server) LoadRouter(ctx context.Context) error {
	list, err := srv.routerModel.GetAll(ctx)
	if err != nil {
		return err
	}
	srv.routerMatcher.Store(router.NewRouter(list))
	return nil
}

func (srv *Server) InitGate() {

}
func (srv *Server) HandleGate(ctx *fasthttp.RequestCtx) {
	srv.gateOnce.Do(func() {
		c, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.LoadRouter(c); err != nil {
			log.Logger.Error("网关服务首次加载路由异常", zap.Error(err))
		}
	})
	routerMatcher, ok := srv.routerMatcher.Load().(*router.RouterMatcher)
	if !ok {
		fasthttputils.RespError(ctx, constant.ErrInitRouterFail)
		return
	}
	host, ok := fasthttputils.GetHeaderHost(ctx)
	if ok {
		ctx.SetUserValue(constant.CtxHeaderHost, host)
	}
	ctx.SetUserValue(constant.CtxRequestPath, fasthttputils.GetRequestPath(ctx))
	ctx.SetUserValue(constant.CtxRequestMethod, string(ctx.Method()))
	ctx.SetUserValue(constant.CtxRequestIP, ctx.RemoteIP().String())
	routerHandler, err := routerMatcher.Match(ctx)
	if err != nil {
		fasthttputils.RespError(ctx, err)
		return
	}
	node, err := routerHandler.GetNode(ctx)
	if err != nil {
		fasthttputils.RespError(ctx, err)
		return
	}
	ctx.SetUserValue(constant.CtxRequestNode, node)
	pass := routerHandler.RequestFilter(ctx)
	if !pass {
		return
	}
	newReq := new(fasthttp.Request)
	ctx.Request.CopyTo(newReq)
	newReq.URI().SetHost(node)
	resp := new(fasthttp.Response)
	if err = srv.fastClient.Do(newReq, resp); err != nil {
		fasthttputils.RespError(ctx, err)
		return
	}
	resp.CopyTo(&ctx.Response)
	routerHandler.ResponseFilter(ctx)
}

func (srv *Server) Start(config config.Config) error {
	ch := make(chan error)
	go func() {
		port := fmt.Sprintf(":%d", config.AdminConfig.Port)
		err := srv.adminEngine.Run(port)
		if err != nil {
			log.Fatal("启动管理后端服务失败", zap.Error(err), zap.String("port", port))
		}
		ch <- err
	}()
	go func() {
		port := fmt.Sprintf(":%d", config.GateConfig.Port)
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", config.GateConfig.Port), srv.gateServer.Handler)
		if err != nil {
			log.Logger.Fatal("启动网关服务失败", zap.Error(err), zap.String("port", port))
		}
		ch <- err
	}()
	return <-ch
}

func NewServer(
	conf config.Config,
	etcdClient *clientv3.Client,
	routerModel dml.RouterModel,
) *Server {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())
	r.Use(middleware.PanicHandle)
	srv := &Server{
		adminEngine:   r,
		routerModel:   routerModel,
		gateOnce:      new(sync.Once),
		routerMatcher: new(atomic.Value),
		etcdClient:    etcdClient,
	}

	route.InitRouter(r.Group(""))
	srv.gateServer = &fasthttp.Server{
		Handler: srv.HandleGate,
	}
	srv.fastClient = &fasthttp.Client{
		ReadTimeout:  time.Duration(conf.GateConfig.Timeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.GateConfig.Timeout) * time.Millisecond,
	}

	return srv
}
