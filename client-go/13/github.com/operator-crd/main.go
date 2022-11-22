package main

import (
	"context"
	"fmt"
	clientset2 "github.com/operator-crd/pkg/generated/clientset/versioned"
	"github.com/operator-crd/pkg/generated/informers/externalversions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {

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

	clientset, err := clientset2.NewForConfig(config)
	if err != nil {
		log.Fatalln("can not create client ")
	}

	list, err := clientset.CrdV1().Foos("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	for _, foo := range list.Items {
		fmt.Println(foo)
	}

	factory := externalversions.NewSharedInformerFactory(clientset, 0)
	factory.Crd().V1().Foos().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// TODO
		},
	})

}
