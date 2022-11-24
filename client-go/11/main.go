package main

import (
	"github.com/darrenli6/CloudNativeGo/client-go/11/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	// 1config

	// 2client

	// 3informer

	// 4 register event

	// 5 start informer

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// 如果在集群内部的话
		inclusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatal("can not get config")
			return
		}
		config = inclusterConfig
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can not create client ")
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()

	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)

	stopCh := make(chan struct{})
	factory.Start(stopCh)
	// 等待同步数据到本地之后
	factory.WaitForCacheSync(stopCh)
	// 同步之后再run
	controller.Run(stopCh)

}
