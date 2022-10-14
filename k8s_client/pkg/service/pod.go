package service

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/client"
	"k8s.io/klog"
)

func GetPod(namespaceName string) ([]v1.Pod, error) {
	ctx := context.Background()

	clientSet, err := client.GetK8SClientSet()

	if err != nil {
		klog.Fatal(err)
		return nil, err
	}
	list, err := clientSet.CoreV1().Pods(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
