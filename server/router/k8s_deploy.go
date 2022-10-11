package router

import (
	v1 "gin-vue-admin/api/v1"
	"github.com/gin-gonic/gin"
)

func InitDeployRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	deployRouter := Router.Group("deploy")
	{
		deployRouter.GET("/list", v1.GetDeployList)
		deployRouter.GET("/:name", v1.GetDeploymentDetail)
		deployRouter.POST("/replicas", v1.IncReplica)
		deployRouter.DELETE("/delete", v1.DeleteDeployment)
	}
	return deployRouter
}
