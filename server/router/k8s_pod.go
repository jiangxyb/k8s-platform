package router

import (
	v1 "gin-vue-admin/api/v1"
	"github.com/gin-gonic/gin"
)

func InitPodRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	podRouter := Router.Group("pod")
	{
		podRouter.GET("/list",v1.GetPods)

	}
	return podRouter
}
