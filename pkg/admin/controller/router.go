package controller

import (
	"ac-gateway/constant"
	"ac-gateway/help/ginutils"
	"ac-gateway/help/strutils"
	"ac-gateway/pkg/admin/form"
	"ac-gateway/pkg/gate/plugins"
	"ac-gateway/pkg/library"
	"ac-gateway/pkg/models/ddl"
	"ac-gateway/pkg/models/dml"
	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/arrays"
	"strings"
	"time"
)

type RouterCtrl struct {
	routerModel dml.RouterModel
}

func NewRouterCtrl(routerModel dml.RouterModel) *RouterCtrl {
	return &RouterCtrl{
		routerModel: routerModel,
	}
}
func (ctrl *RouterCtrl) CheckRouter(router ddl.Router) error {
	if len(router.Upstream.Nodes) == 0 {
		return constant.ErrRouterEmptyNodes
	}
	if len(router.Methods) != 0 {
		for _, m := range router.Methods {
			m = strings.ToUpper(strings.TrimSpace(m))
			if arrays.ContainsString(constant.DefaultHttpMethod, m) == -1 {
				return constant.ErrHttpMethod
			}
		}
	}
	for node := range router.Upstream.Nodes {
		if strutils.IsEmpty(node) {
			return constant.ErrRouterNodeSpace
		}
	}
	router.Upstream.Type = strings.TrimSpace(router.Upstream.Type)
	if router.Upstream.Type != constant.RouteUpstreamTypeRoundRobin &&
		router.Upstream.Type != constant.RouteUpstreamTypeIPHash &&
		strutils.IsEmpty(router.Upstream.Type) {
		return constant.ErrRouterUpstreamType
	}
	if len(router.Plugins.ExtPluginPreReq.Conf) != 0 {
		for _, c := range router.Plugins.ExtPluginPreReq.Conf {
			if strutils.IsEmpty(c.Name) {
				return constant.ErrRouterPluginEmpty
			}
			_, ok := plugins.GetPlugin(c.Name)
			if !ok {
				return constant.ErrRouterNotFoundPlugin
			}
		}
	}
	if len(router.Plugins.ExtPluginPostResp.Conf) != 0 {
		for _, c := range router.Plugins.ExtPluginPreReq.Conf {
			if strutils.IsEmpty(c.Name) {
				return constant.ErrRouterPluginEmpty
			}
			_, ok := plugins.GetPlugin(c.Name)
			if !ok {
				return constant.ErrRouterNotFoundPlugin
			}
		}
	}
	return nil
}
func (ctrl *RouterCtrl) Add(c *gin.Context) {
	requestForm := form.RouterAddForm{}
	if err := c.Bind(&requestForm); err != nil {
		return
	}
	router := requestForm.ToRouter()
	if err := ctrl.CheckRouter(router); err != nil {
		ginutils.RespError(c, err)
		return
	}
	router.Id = time.Now().UnixMilli()
	if strutils.IsNotEmpty(router.Name) {
		list, err := ctrl.routerModel.GetAll(c)
		if err != nil {
			ginutils.RespError(c, err)
			return
		}
		for _, r := range list {
			if r.Name == router.Name {
				ginutils.RespError(c, constant.ErrRouterNameExists)
				return
			}
		}
	}

	router, err := ctrl.routerModel.Save(c, router)
	if err != nil {
		ginutils.RespError(c, err)
		return
	}
	ginutils.RespSuccessWithData(c, router)
}

func (ctrl *RouterCtrl) Get(c *gin.Context) {
	requestForm := form.IdFrom{}
	if err := c.Bind(&requestForm); err != nil {
		return
	}
	router, exists, err := ctrl.routerModel.GetById(c, requestForm.Id)
	if err != nil {
		ginutils.RespError(c, err)
		return
	}
	if !exists {
		ginutils.RespError(c, constant.ErrRouterNotFound)
		return
	}
	ginutils.RespSuccessWithData(c, router)
}

func (ctrl *RouterCtrl) GetAll(c *gin.Context) {
	list, err := ctrl.routerModel.GetAll(c)
	if err != nil {
		ginutils.RespError(c, err)
		return
	}
	ginutils.RespSuccessWithData(c, list)
}

func (ctrl *RouterCtrl) Update(c *gin.Context) {
	requestForm := form.RouterUpdateForm{}
	if err := c.Bind(&requestForm); err != nil {
		return
	}
	router := requestForm.ToRouter()
	if err := ctrl.CheckRouter(router); err != nil {
		ginutils.RespError(c, err)
		return
	}
	list, err := ctrl.routerModel.GetAll(c)
	if err != nil {
		ginutils.RespError(c, err)
		return
	}
	exists := false
	for _, r := range list {
		if r.Id == router.Id {
			exists = true
		}
		if r.Id != router.Id && r.Name == router.Name {
			ginutils.RespError(c, constant.ErrRouterNameExists)
			return
		}
	}

	if !exists {
		ginutils.RespError(c, constant.ErrRouterNotFound)
		return
	}
	list = library.FormatRouterList(list)
	if router, err = ctrl.routerModel.Save(c, router); err != nil {
		ginutils.RespError(c, err)
		return
	}
	ginutils.RespSuccessWithData(c, router)
}

func (ctrl *RouterCtrl) Delete(c *gin.Context) {
	requestForm := form.IdFrom{}
	if err := c.Bind(&requestForm); err != nil {
		return
	}
	if err := ctrl.routerModel.Delete(c, requestForm.Id); err != nil {
		ginutils.RespError(c, err)
		return
	}
	ginutils.RespSuccess(c)
}
