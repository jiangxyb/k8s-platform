package v1

import (
	"context"
	"errors"
	"fmt"
	"gin-vue-admin/global"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func GetDeployment(ns string, name string) *v1.Deployment {
	ctx := context.Background()
	getopt := metav1.GetOptions{}
	depDetail, _ := global.ClientSet.AppsV1().Deployments(ns).Get(ctx, name, getopt)
	return depDetail
}

func GetDeploymentByWatch(ns string, name string) (*v1.Deployment, error) {
	if deploys, ok := global.DeployMap.Load(ns); ok {
		for _, deploy := range deploys.([]*v1.Deployment) {
			if deploy.Name == name {
				return deploy, nil
			}
		}
	}
	return nil, fmt.Errorf("record not found")
}

func DeploysByNS(ns string) ([]*v1.Deployment, error) {
	if list, ok := global.DeployMap.Load(ns); ok {
		return list.([]*v1.Deployment), nil
	}
	return nil, fmt.Errorf("record not found")
}

func GetPodsByDeploy(ns string, deploy *v1.Deployment) []*corev1.Pod {
	rs, _ := GetRSByDeploymentWithWatch(ns, deploy)
	labels, _ := metav1.LabelSelectorAsMap(rs.Spec.Selector)
	fmt.Println(labels)
	podList, _ := global.PodMap.ListByLabels(ns, labels)
	fmt.Println(len(podList))
	return podList
}

func GetLabelSelectorByDeploy(deploy *v1.Deployment) string {
	str := ""
	for k, v := range deploy.Spec.Selector.MatchLabels {
		if str != "" {
			str += ","
		}
		str += fmt.Sprintf("%s=%s", k, v)
	}
	fmt.Println("字符串是：", str)
	return str
}

func GetImagesByContainers(containers []corev1.Container) []string {
	images := make([]string, len(containers))
	for i, v := range containers {
		images[i] = v.Image
	}
	return images
}

func GetRSListByNamespace(ns string) []*v1.ReplicaSet {
	opts := metav1.ListOptions{}
	RSList, _ := global.ClientSet.AppsV1().ReplicaSets(ns).List(context.Background(), opts)
	rSs := make([]*v1.ReplicaSet, len(RSList.Items))
	for i, rs := range RSList.Items {
		rSs[i] = &rs
		fmt.Println(rs.Name)
	}
	return rSs
}

// ns下要有这个deployment才能获取到rs
func GetRSByDeployment(ns string, dep *v1.Deployment) (*v1.ReplicaSet, error) {
	selector, _ := metav1.LabelSelectorAsSelector(dep.Spec.Selector)

	// 先用label匹配rs
	opts := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	RSList, _ := global.ClientSet.AppsV1().ReplicaSets(ns).List(context.Background(), opts)
	// 再用注解和依赖资源筛选
	for _, rs := range RSList.Items {
		if IsCurrentRsByDep(dep, &rs) {
			return &rs, nil
		}
	}
	return nil, errors.New("not found")
}

// ns下要有这个deployment才能获取到rs
func GetRSByDeploymentWithWatch(ns string, dep *v1.Deployment) (*v1.ReplicaSet, error) {
	RSList, err := global.RsMap.ListByNameSpace(ns)
	if err != nil {
		log.Println(err)
	}
	// 用注解和依赖资源筛选
	for _, rs := range RSList {
		if IsCurrentRsByDep(dep, rs) {
			return rs, nil
		}
	}
	return nil, errors.New("not found")
}

func GetRsLabelByDeployment(dep *v1.Deployment, rslist []*v1.ReplicaSet) (map[string]string, error) {
	for _, item := range rslist {
		if IsCurrentRsByDep(dep, item) {
			s, err := metav1.LabelSelectorAsMap(item.Spec.Selector)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, nil
}

// 用注解和依赖资源筛选
func IsCurrentRsByDep(dep *v1.Deployment, rs *v1.ReplicaSet) bool {
	if rs.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != dep.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] {
		return false
	}
	for _, ref := range rs.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == dep.Name {
			return true
		}
	}
	return false
}
