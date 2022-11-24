package pkg

import (
	"context"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"reflect"
	"time"
)

const (
	workNum  = 5
	maxRetry = 10
)

type controller struct {
	// 操作 service ingress
	client kubernetes.Interface
	// 获取资源对象的状态，避免跟api server交互
	ingressLister v1.IngressLister
	serviceLister coreLister.ServiceLister
	// 限速队列
	queue workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, serviceInfomer informer.ServiceInformer, ingressInformer netInformer.IngressInformer) controller {

	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInfomer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}
	serviceInfomer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return c

}

// 都会用到queue
func (c *controller) addService(obj interface{}) {
	// 新建一个对象
	c.enqueue(obj)

}

func (c *controller) updateService(obj interface{}, newobj interface{}) {
	//todo 比较annotation
	// 对象一致就不进行处理
	if reflect.DeepEqual(obj, newobj) {
		return
	}
	// 不一致的进队列
	c.enqueue(newobj)

}

func (c *controller) enqueue(obj interface{}) {
	//通过object拿到key
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	c.queue.Add(key)
}

func (c *controller) deleteIngress(obj interface{}) {

	ingress := obj.(*v12.Ingress)

	ownerReference := v13.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	// 确定是否是service
	if ownerReference.Kind != "Service" {
		return
	}
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) Run(stopCh chan struct{}) {
	// 消费数据
	for i := 0; i < workNum; i++ {
		// 队列
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	// 关闭
	<-stopCh
}

func (c *controller) worker() {
	// 不停地获取key 处理循环
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {

	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	// 处理key 完成之后移除
	defer c.queue.Done(item)
	key := item.(string)

	// 重试 错误处理
	err := c.syncService(key)
	if err != nil {
		c.handleError(key, err)
	}
	return true
}

func (c *controller) syncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespaceKey).Get(name)

	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]

	ingress, err := c.ingressLister.Ingresses(namespaceKey).Get(name)
	if err != nil {
		return err
	}
	// 有的 ingress/http的话
	if ok && errors.IsNotFound(err) {
		// 创建ingress
		ig := c.contructIngress(namespaceKey, name)
		_, err := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v13.CreateOptions{})
		if err != nil {
			return nil
		}
	} else if !ok && ingress != nil {
		// 删除ingress
		err := c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, v13.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) handleError(key string, err error) {

	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}
	runtime.HandleError(err)
	// 清空key
	c.queue.Forget(key)
}

func (c *controller) contructIngress(namespacekey string, name string) *v12.Ingress {

	ingress := v12.Ingress{}
	ingress.Name = name
	ingress.Namespace = namespacekey
	pathType := v12.PathTypePrefix
	ingress.Spec = v12.IngressSpec{
		Rules: []v12.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: v12.IngressRuleValue{
					HTTP: &v12.HTTPIngressRuleValue{
						Paths: []v12.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: v12.IngressBackend{
									Service: &v12.IngressServiceBackend{
										Name: name,
										Port: v12.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &ingress
}
