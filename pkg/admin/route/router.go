package route

import (
	"ac-gateway/pkg/admin/controller"
	"ac-gateway/pkg/di"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(router *gin.RouterGroup) {
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "", []byte("You Build It, You Run It."))
	})
	adminRouter := router.Group("/admin/router")
	di.MustInvoke(func(routerCtrl *controller.RouterCtrl) {
		adminRouter.GET("/get", routerCtrl.Get)
		adminRouter.POST("/add", routerCtrl.Add)
		adminRouter.POST("/update", routerCtrl.Update)
		adminRouter.POST("/delete", routerCtrl.Delete)
		adminRouter.GET("/get_all", routerCtrl.GetAll)
	})

}
