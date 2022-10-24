package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
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

	for c.processItme() {

	}
}

func (c *controller) processItme() bool {

	item, shutdown := c.queue.Get() // 从队列中获取项目

	if shutdown {
		return false
	}
	defer c.queue.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)

	if err != nil {
		fmt.Printf("getting key form cache %s", err.Error())
	}
	//引入名称和命名空间
	ns, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		fmt.Printf("slliting key into namespace and name %s \n", err.Error())
		return false
	}

	// 现在将特定的控制器创建serviec,首先调用deployment
	err = c.syncDeployment(ns, name)
	if err != nil {
		fmt.Printf("syncing deployment %s \n", err.Error())
		return false
	}

	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	dep, err := c.depLister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Printf(" getting deployment from listering %s \n", err.Error())

	}
	// 创建service
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: depLabels(*dep),
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     80,
					NodePort: 30080,
				},
			},
		},
	}

	_, err = c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("create service %s \n", err.Error())
	}
	return nil
}

func depLabels(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("注册的添加函数调用了，创建deployment")
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("注册的删除函数调用了,删除deployment")
	c.queue.Add(obj)
}
