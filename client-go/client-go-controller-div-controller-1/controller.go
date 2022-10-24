package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset      kubernetes.Interface
	depLister      appslisters.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:      clientset,                        // client来初始化
		depLister:      depInformer.Lister(),             // 将是部署列出
		depCacheSynced: depInformer.Informer().HasSynced, // 注册缓冲同步信息
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kubedarren-expose"),
	}

	// 一旦有个部署通知器 讨论注册函数
	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			DeleteFunc: handleDel,
		},
	)

	return c
}

func (c *controller) run(ch <-chan struct{}) {

	fmt.Println("starting controller")
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("waiting for cache to be synced\n")

	}

	go wait.Until(c.worker, time.Second, ch)
	//等待直到它的作用是它在每个持续时间后调用一个特定函数，直到该特定channel关闭，所以如果
	//这个函数运行我们要指定的函数，这个函数将在每个周期之后运行，直到关闭
	//所以如果我们不关闭这个channel,我们将通过这个函数每次我们指定之后都会贝调用，所以我们
	//称它为调用函数c.worker,指定时间秒，这样我们指定的通道就会将在这里运行，这就是一个很好的例程
	<-ch
}

func (c *controller) worker() {

}

func handleAdd(obj interface{}) {
	fmt.Println("注册的添加函数调用了，创建deployment")
}

func handleDel(obj interface{}) {
	fmt.Println("注册的删除函数调用了,删除deployment")
}
