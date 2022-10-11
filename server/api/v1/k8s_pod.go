package v1

import (
	"context"
	"gin-vue-admin/global"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetPods(ctx *gin.Context) {
	c := context.Background()
	opt := metav1.ListOptions{}
	podList, _ := global.ClientSet.CoreV1().Pods("default").List(c, opt)
	ctx.JSON(http.StatusOK, &podList)
}

func GetPodJson(ctx *gin.Context) {
	ns := ctx.DefaultQuery("ns", "default")
	podName := ctx.DefaultQuery("pod", "default")
	pod := global.PodMap.Get(ns, podName)
	ctx.JSON(http.StatusOK, &pod)
}
