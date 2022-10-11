package global

import (
	"fmt"
	"gorm.io/gorm"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sync"
	"time"
)

type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

// 保存Pod集合
type PodMapStruct struct {
	data sync.Map // [key string] []*v1.Pod    key=>namespace
}

func (this *PodMapStruct) Add(pod *corev1.Pod) {
	if list, ok := this.data.Load(pod.Namespace); ok {
		list = append(list.([]*corev1.Pod), pod)
		this.data.Store(pod.Namespace, list)
	} else {
		this.data.Store(pod.Namespace, []*corev1.Pod{pod})
	}
}
func (this *PodMapStruct) Update(pod *corev1.Pod) error {
	if list, ok := this.data.Load(pod.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Pod) {
			if range_pod.Name == pod.Name {
				list.([]*corev1.Pod)[i] = pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found", pod.Name)
}
func (this *PodMapStruct) Delete(pod *corev1.Pod) {
	if list, ok := this.data.Load(pod.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Pod) {
			if range_pod.Name == pod.Name {
				newList := append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				this.data.Store(pod.Namespace, newList)
				break
			}
		}
	}
}

func (this *PodMapStruct) Get(ns string, podName string) *corev1.Pod {
	if list, ok := this.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			if pod.Name == podName {
				return pod
			}
		}
	}
	return nil
}

//根据多个rs的标签获取 POD列表
func (this *PodMapStruct) ListByLabels(ns string, labels []map[string]string) ([]*corev1.Pod, error) {
	ret := make([]*corev1.Pod, 0)
	if list, ok := this.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			for _, l := range labels {
				isMatchedPod := true
				for k, v := range l {
					if value, ok := pod.Labels[k]; !ok || value != v {
						isMatchedPod = false
						break
					}
				}
				if isMatchedPod {
					ret = append(ret, pod)
				}

			}
		}
		return ret, nil
	}
	return nil, fmt.Errorf("pods not found ")
}

// ReplicaSet 集合
type RsMapStruct struct {
	data sync.Map // [key string] []*appv1.ReplicaSet    key=>namespace
}

func (this *RsMapStruct) Add(rs *appv1.ReplicaSet) {
	if list, ok := this.data.Load(rs.Namespace); ok {
		list = append(list.([]*appv1.ReplicaSet), rs)
		this.data.Store(rs.Namespace, list)
	} else {
		this.data.Store(rs.Namespace, []*appv1.ReplicaSet{rs})
	}
}
func (this *RsMapStruct) Update(rs *appv1.ReplicaSet) error {
	if list, ok := this.data.Load(rs.Namespace); ok {
		for i, range_rs := range list.([]*appv1.ReplicaSet) {
			if range_rs.Name == rs.Name {
				list.([]*appv1.ReplicaSet)[i] = rs
			}
		}
		return nil
	}
	return fmt.Errorf("rs-%s not found", rs.Name)
}
func (this *RsMapStruct) Delete(rs *appv1.ReplicaSet) {
	if list, ok := this.data.Load(rs.Namespace); ok {
		for i, range_rs := range list.([]*appv1.ReplicaSet) {
			if range_rs.Name == rs.Name {
				newList := append(list.([]*appv1.ReplicaSet)[:i], list.([]*appv1.ReplicaSet)[i+1:]...)
				this.data.Store(rs.Namespace, newList)
				break
			}
		}
	}
}

//普普通通的函数， 就是根据ns获取 对应的rs列表
func (this *RsMapStruct) ListByNameSpace(ns string) ([]*appv1.ReplicaSet, error) {
	if list, ok := this.data.Load(ns); ok {
		return list.([]*appv1.ReplicaSet), nil
	}
	return nil, fmt.Errorf("pods not found ")
}

type EventMapStruct struct {
	data sync.Map // value=> *v1.Event
	// key=>namespace+"_"+kind+"_"+name 这里的name 不一定是pod
}

func (this *EventMapStruct) GetMessage(ns string, kind string, name string) string {
	key := fmt.Sprintf("%s_%s_%s", ns, kind, name)
	if v, ok := this.data.Load(key); ok {
		return v.(*corev1.Event).Message
	}

	return ""
}

func (this *EventMapStruct) Store(key string, event *corev1.Event) {
	this.data.Store(key, event)
}

func (this *EventMapStruct) Delete(key string) {
	this.data.Delete(key)
}
