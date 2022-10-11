package initialize

import (
	"fmt"
	"gin-vue-admin/global"
	corev1 "k8s.io/api/core/v1"
)

func events() {
	podInformer := fact.Core().V1().Events()
	podInformer.Informer().AddEventHandler(&EventHandler{})
}

type EventHandler struct{}

func (this *EventHandler) storeData(obj interface{}, isdelete bool) {
	if event, ok := obj.(*corev1.Event); ok {
		key := fmt.Sprintf("%s_%s_%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
		if !isdelete {
			global.EventMap.Store(key, event)
		} else {
			global.EventMap.Delete(key)
		}
	}
}
func (this *EventHandler) OnAdd(obj interface{}) {
	this.storeData(obj, false)
}
func (this *EventHandler) OnUpdate(oldObj, newObj interface{}) {
	this.storeData(newObj, false)
}
func (this *EventHandler) OnDelete(obj interface{}) {
	this.storeData(obj, true)
}
