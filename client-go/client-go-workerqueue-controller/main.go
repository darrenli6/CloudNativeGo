package main

import (
	"flag"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	// 等待 知道工作队列中有一个新项
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// 告诉队列我们已经处理完这个键，这允许安全的并行处理 因为两个具有相同的秘钥的pod永远不会并行处理
	defer c.queue.Done(key)

	// 调用包含业务逻辑的方法
	err := c.syncToStdout(key.(string))

	// 如果执行业务逻辑出现错误 处理错误
	c.handleErr(err, key)

	return true

}

// 处理错误
func (c *Controller) handleErr(err error, key interface{}) {

	if err == nil {
		// 这个确保了以后对这个键的更新处理 不会因为过时的错误历史而延迟
		c.queue.Forget(key)
		return
	}
	// 如果出现问题，这个控制器会重试5次，在那以后，它就会停止尝试
	if c.queue.NumRequeues(key) < 5 {

		klog.Infof(" error syncing pod %v %v", key, err)
		// 重新排队
		c.queue.AddRateLimited(key)
		return
	}
	c.queue.Forget(key)
	// 向外部报告
	runtime.HandleError(err)
	klog.Infof(" Dropping pod %q out of the queue %v", key, err)
}

// syncToStdout是控制器的业务逻辑 在这个控制器中，它只是将关于pod的信息打印到stdout
// 在发生错误的情况下， 它必须简单地返回错误
// 重试逻辑 不应该是业务逻辑的一部分
func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		// 暖化我的们的缓冲
		fmt.Printf("pod %s does not exist anymore \n", key)
	} else {
		// 注意 如果你有一个本地控制资源，你也必须检查uid 这是依赖实际的实例 检测一个pod  被创建了相同的名称
		fmt.Printf(" sync/add/update for pod %s \n", obj.(*v1.Pod).GetName())

	}
	return nil

}

// Run开始观察和同步
func (c *Controller) Run(workers int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// 我们完工之后让任务停下来
	defer c.queue.ShutDown()

	klog.Info("Staring pod controller ")

	go c.informer.Run(stopCh)

	// 在开始处理队列中的项目之前，等待所有设计的缓冲被同步
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Time out waiting for caches to sync"))
		return
	}
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	klog.Info(" stopping pod controller ")

}

func (c *Controller) runWorker() {
	for c.processNextItem() {

	}
}

func main() {

	var kubeconfig string
	var master string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file ")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()

	// 创建链接
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())

	// 创建workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// 在informer的帮助下绑定工作队列到缓冲， 我们可以确保每当缓冲被更新的时候，pod键就添加到工作队列中
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{

		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceIndexFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	indexer.Add(&v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "test-node-local-dns",
			Namespace: v1.NamespaceDefault,
		},
	})

	stop := make(chan struct{})

	defer close(stop)

	go controller.Run(1, stop)

	select {}

}
