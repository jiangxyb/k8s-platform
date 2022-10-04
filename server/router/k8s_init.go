package router

import "github.com/gin-gonic/gin"

func InitKubernetesRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	K8sRouter := Router.Group("kubernetes")
	{
		InitDeployRouter(K8sRouter)
		InitPodRouter(K8sRouter)
	}
	return K8sRouter
}
