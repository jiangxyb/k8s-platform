package v1

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetDeployList(ctx *gin.Context) {
	deploys, _ := DeploysByNS("default")
	deployList := model.Deployments{
		Items: deploys,
	}

	ctx.JSON(http.StatusOK, deployList)
}

func GetDeploymentDetail(ctx *gin.Context) {
	name := ctx.Param("name")
	fmt.Println(name)
	detail, _ := GetDeploymentByWatch("default", name)

	iNames := make([]string, 0)
	for _, v := range detail.Spec.Template.Spec.Containers {
		iNames = append(iNames, v.Image)
	}
	podList := GetPodsByDeploy("default", detail)
	pods := make([]*model.Pod, len(podList))
	msg := ""
	for i, pod := range podList {
		for _, condition := range pod.Status.Conditions {
			if condition.Status != "True" {
				msg = fmt.Sprintf("%v,%v", condition.Reason, condition.Message)
				break
			}
		}
		pods[i] = &model.Pod{
			Name:         pod.Name,
			Ns:           pod.Namespace,
			Images:       GetImagesByContainers(pod.Spec.Containers),
			NodeName:     pod.Spec.NodeName,
			CreateTime:   pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Phase:        string(pod.Status.Phase),
			Message:      msg,
			EventMsg:     global.EventMap.GetMessage("default", "Pod", pod.Name),
			RestartCount: pod.Status.ContainerStatuses[0].RestartCount,
			Ready:        GetPodIsReady(pod),
		}
	}

	replicas := make([]int32, 3)
	replicas[0] = *detail.Spec.Replicas
	replicas[1] = detail.Status.AvailableReplicas
	replicas[2] = detail.Status.UnavailableReplicas
	deployDetail := model.DeployDetail{
		Name:       detail.Name,
		ImageNames: iNames,
		Replicas:   replicas,
		Pods:       pods,
	}
	ctx.JSON(http.StatusOK, &deployDetail)
}

func IncReplica(ctx *gin.Context) {
	req := model.ReqReplica{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		panic(err)
	}
	opt := metav1.GetOptions{}
	scale, _ := global.ClientSet.AppsV1().Deployments(req.NameSpace).GetScale(context.Background(), req.Deploy, opt)
	if req.Dec {
		scale.Spec.Replicas -= 1
	} else {
		scale.Spec.Replicas += 1
	}
	opts := metav1.UpdateOptions{}
	_, err = global.ClientSet.AppsV1().Deployments(req.NameSpace).UpdateScale(context.Background(), req.Deploy, scale, opts)
	if err != nil {
		panic(err)
	}
	response.Ok(ctx)
}

func DeleteDeployment(ctx *gin.Context) {
	ns := ctx.DefaultQuery("ns", "default")
	deploy := ctx.Query("deploy")
	err := global.ClientSet.AppsV1().Deployments(ns).Delete(context.Background(), deploy, metav1.DeleteOptions{})
	if err != nil {
		global.GVA_LOG.Error("删除deployment失败!", zap.Any("err", err))
		response.FailWithMessage("删除deployment失败", ctx)
		return
	}
	response.Ok(ctx)
}
