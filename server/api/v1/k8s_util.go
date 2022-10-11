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
	labels, _ := GetRsListLabel(rs)
	fmt.Println(labels)
	podList, _ := global.PodMap.ListByLabels(ns, labels)
	fmt.Println(len(podList))
	return podList
}

func GetRsListLabel(rslist []*v1.ReplicaSet) ([]map[string]string, error) {

	ret := make([]map[string]string, 0)
	for _, item := range rslist {
		s, err := metav1.LabelSelectorAsMap(item.Spec.Selector)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
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
func GetRSByDeploymentWithWatch(ns string, dep *v1.Deployment) ([]*v1.ReplicaSet, error) {
	RSList, err := global.RsMap.ListByNameSpace(ns)
	if err != nil {
		log.Println(err)
	}

	ret := make([]*v1.ReplicaSet, 0)
	for _, rs := range RSList {
		if IsRsOfDep(dep, rs) {
			ret = append(ret, rs)
		}
	}
	if len(ret) != 0 {
		return ret, nil
	}
	return nil, errors.New("not found")
}

func GetRsLabelByDeployment(dep *v1.Deployment, rslist []*v1.ReplicaSet) (map[string]string, error) {
	for _, item := range rslist {
		if IsRsOfDep(dep, item) {
			s, err := metav1.LabelSelectorAsMap(item.Spec.Selector)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, nil
}

//判断 rs 是否属于 某个 dep
func IsRsOfDep(dep *v1.Deployment, set *v1.ReplicaSet) bool {
	for _, ref := range set.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == dep.Name {
			return true
		}
	}
	return false
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

//判断POD是否就绪
func GetPodIsReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "ContainersReady" && condition.Status != "True" {
			return false
		}
	}
	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}
	return true
}
