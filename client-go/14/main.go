package main

import (
	"context"
	v1 "darren.tech/pkg/apis/darren.tech/v1"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalln(err)
	}

	config.APIPath = "/apis/"
	config.NegotiatedSerializer = v1.Codecs.WithoutConversion()
	config.GroupVersion = &v1.GroupVersion

	client, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatalln(err)
	}
	foo := &v1.Foo{}

	err = client.Get().Namespace("default").Resource("foos").Name("crd-test").Do(context.TODO()).Into(foo)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(foo.Spec)

	newobj := foo.DeepCopy()
	newobj.Spec.Name = "ahhaha"

	fmt.Println(newobj.Spec)

}
