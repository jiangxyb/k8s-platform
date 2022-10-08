package initialize

import (
	"gin-vue-admin/global"
	appv1 "k8s.io/api/apps/v1"
	"log"
)

func replicaset() {
	rsInformer := fact.Apps().V1().ReplicaSets()
	rsInformer.Informer().AddEventHandler(&RsHandler{})
}

type RsHandler struct{}

func (this *RsHandler) OnAdd(obj interface{}) {
	global.RsMap.Add(obj.(*appv1.ReplicaSet))
}
func (this *RsHandler) OnUpdate(oldObj, newObj interface{}) {
	err := global.RsMap.Update(newObj.(*appv1.ReplicaSet))
	if err != nil {
		log.Println(err)
	}
}
func (this *RsHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*appv1.ReplicaSet); ok {
		global.RsMap.Delete(d)
	}
}
