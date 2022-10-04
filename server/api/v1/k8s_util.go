package v1

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployment(ns string,name string ) *v1.Deployment  {
	ctx:=context.Background()
	getopt:=metav1.GetOptions{}
	depDetail,_:=global.ClientSet.AppsV1().Deployments(ns).Get(ctx,name,getopt)
	return depDetail
}

func GetPodsByDeploy(ns string,deploy *v1.Deployment) *corev1.PodList {
	c := context.Background()
	opt := metav1.ListOptions{
		LabelSelector: GetLabelSelectorByDeploy(deploy),
	}
	podList,_ := global.ClientSet.CoreV1().Pods(ns).List(c,opt)
	return podList
}

func GetLabelSelectorByDeploy(deploy *v1.Deployment) string {
	str := ""
	for k,v := range deploy.Spec.Selector.MatchLabels {
		if str != "" {
			str+= ","
		}
		str += fmt.Sprintf("%s=%s",k,v)
	}
	fmt.Println("字符串是：",str)
	return str
}

func GetImagesByContainers(containers []corev1.Container) []string {
	images := make([]string,len(containers))
	for i,v := range containers {
		images[i] = v.Image
	}
	return images
}

func GetImagesByDeploy() {

}