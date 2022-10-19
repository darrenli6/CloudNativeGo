package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	kubeconfig := flag.String("kubeconfig", homedir.HomeDir()+"/.kube/config", "location to your kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Printf(" error %s building config from flags \n", err.Error())
		config, err = rest.InClusterConfig()

		if err != nil {
			fmt.Printf("err %s getting in cluster", err.Error())
		}
	}

	// 实例化 client set

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Printf(" error %s ", err.Error())
	}

	ctx := context.Background()
	fmt.Println("获取default namespace 下的pod")

	// 获取default的pod名字
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Printf(" get pod happened error  %s", err.Error())

	}

	for _, pod := range pods.Items {
		fmt.Printf("%s \n", pod.Name)
	}

	fmt.Println("获取default namespace下的deployment的名字")

	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Printf(" get deployment happened error  %s", err.Error())

	}
	for _, d := range deployments.Items {
		fmt.Printf("%s\n", d.Name)
	}

	fmt.Println("获取kube-system namespace下的daemonset的名字 ")
	//3、获取kube-system下daemonset的资源名字
	daemonsets, err := clientset.AppsV1().DaemonSets("kube-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("listing daemonsets %s\n", err.Error())
	}
	for _, ds := range daemonsets.Items {
		fmt.Printf("%s\n", ds.Name)
	}

	fmt.Println("获取get node的方法 ")
	//4、获取get node的的名字
	node, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("listing Node %s\n", err.Error())
	}
	for _, no := range node.Items {
		fmt.Printf("%s\n", no.Name)
	}

	fmt.Println("获取kube-system svc")
	//5、获取kube-system下的service
	svc, err := clientset.CoreV1().Services("kube-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("listing service %s\n", err.Error())
	}
	for _, service := range svc.Items {
		fmt.Printf("%s\n", service.Name)
	}

}
