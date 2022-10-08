package initialize

import (
	"fmt"
	"gin-vue-admin/global"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var fact informers.SharedInformerFactory

func Deployment() {
	fact = informers.NewSharedInformerFactory(global.ClientSet, 0)
	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(&DepHandler{})

	pod()
	replicaset()

	fact.Start(wait.NeverStop)
}

// 回调类
type DepHandler struct{}

// 资源增加时回调
func (this *DepHandler) OnAdd(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {

		if deploys, ok := global.DeployMap.Load(dep.Namespace); ok {
			for i, deploy := range deploys.([]*v1.Deployment) {
				if deploy.Name == dep.Name {
					deploys.([]*v1.Deployment)[i] = deploy
					return
				}
			}
			global.DeployMap.Store(dep.Namespace, append(deploys.([]*v1.Deployment), dep))
		} else {
			global.DeployMap.Store(dep.Namespace, []*v1.Deployment{dep})
		}
	}
}

// 资源内部变动时回调
func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	if dep, ok := newObj.(*v1.Deployment); ok {
		if deploys, ok := global.DeployMap.Load(dep.Namespace); ok {
			for i, deploy := range deploys.([]*v1.Deployment) {
				if deploy.Name == dep.Name {
					deploys.([]*v1.Deployment)[i] = dep
				}
			}
			fmt.Println()
			return
		}
	}
}

// 资源被删除时回调
func (this *DepHandler) OnDelete(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		if deploys, ok := global.DeployMap.Load(dep.Namespace); ok {
			for i, deploy := range deploys.([]*v1.Deployment) {
				if deploy.Name == dep.Name {
					newDeploys := append(deploys.([]*v1.Deployment)[:i], deploys.([]*v1.Deployment)[i+1:]...)
					global.DeployMap.Store(dep.Namespace, newDeploys)
					return
				}
			}
		}
	}
}

func listAndWatch() {
	_, c := cache.NewInformer(cache.NewListWatchFromClient(global.ClientSet.AppsV1().RESTClient(),
		"deployments", "default", fields.Everything()),
		&v1.Deployment{},
		0,
		&DepHandler{},
	)
	c.Run(wait.NeverStop)
	select {}
}
