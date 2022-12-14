package global

import (
	"gin-vue-admin/utils/timer"
	"k8s.io/client-go/kubernetes"
	"sync"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"gin-vue-admin/config"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_DB     *gorm.DB
	GVA_REDIS  *redis.Client
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	//GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}
	ClientSet               *kubernetes.Clientset
	// key: name of namespace,value: []*v1.Deployment
	DeployMap sync.Map
	PodMap    PodMapStruct //作为全局对象
	RsMap     RsMapStruct
	EventMap  EventMapStruct
)
