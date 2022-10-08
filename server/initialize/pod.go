package initialize

import (
	"gin-vue-admin/global"
	corev1 "k8s.io/api/core/v1"
	"log"
)

func pod() {
	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&PodHandler{})
}

type PodHandler struct{}

func (this *PodHandler) OnAdd(obj interface{}) {
	global.PodMap.Add(obj.(*corev1.Pod))
}
func (this *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	err := global.PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		log.Println(err)
	}
}
func (this *PodHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*corev1.Pod); ok {
		global.PodMap.Delete(d)
	}
}
