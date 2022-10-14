package router

import (
	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/apis"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", apis.Ping)
	r.GET("/namespace", apis.GetNamespace)
	// r.GET("/pod", apis.GetPods)
	r.GET("/namespace/:namespaceName/pods", apis.GetPods)
}
