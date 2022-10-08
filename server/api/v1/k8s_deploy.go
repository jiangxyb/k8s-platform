package v1

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetDeployList(ctx *gin.Context) {
	deploys, _ := DeploysByNS("default")
	deployList := Deployments{
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
	pods := make([]*Pod, len(podList))
	for i, item := range podList {
		pods[i] = &Pod{
			Name:       item.Name,
			Images:     GetImagesByContainers(item.Spec.Containers),
			NodeName:   item.Spec.NodeName,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}

	replicas := make([]int32, 3)
	replicas[0] = *detail.Spec.Replicas
	replicas[1] = detail.Status.AvailableReplicas
	replicas[2] = detail.Status.UnavailableReplicas
	deployDetail := DeployDetail{
		Name:       detail.Name,
		ImageNames: iNames,
		Replicas:   replicas,
		Pods:       pods,
	}
	ctx.JSON(http.StatusOK, &deployDetail)
}

func IncReplica(ctx *gin.Context) {
	req := ReqReplica{}
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

type Deployments struct {
	Items []*v1.Deployment `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ReqReplica struct {
	NameSpace string `json:"ns"`
	Deploy    string `json:"deploy"`
	Dec       bool   `json:"dec"`
}

type DeployDetail struct {
	Name       string   `json:"name"`
	ImageNames []string `json:"images_name"`
	Replicas   []int32  `json:"replicas"` // index表示表示不同的含义，0 总，1 可用，2 不可用
	Pods       []*Pod   `json:"pods"`
}

type Pod struct {
	Name       string   `json:"name"`
	Images     []string `json:"images"`
	NodeName   string   `json:"node_name"`
	CreateTime string   `json:"create_time"`
}
