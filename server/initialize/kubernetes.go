package initialize

import (
	"gin-vue-admin/global"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func Kubernetes() {
	config:=&rest.Config{
		Host: global.GVA_CONFIG.Kubernetes.IP,
		BearerToken: global.GVA_CONFIG.Kubernetes.Token,
	}
	config.Insecure = true
	var err error
	global.ClientSet,err =kubernetes.NewForConfig(config)
	if err!=nil{
		log.Fatal(err)
	}
}
