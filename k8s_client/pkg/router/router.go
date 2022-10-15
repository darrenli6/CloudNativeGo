package router

import (
	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/apis"
	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	middleware.InitMiddleware(r)
	r.GET("/ping", apis.Ping)
	r.GET("/namespace", apis.GetNamespace)

	r.GET("/namespace/:namespaceName/pods", apis.GetPods)
	r.GET("/namespace/:namespaceName/pod/:podName/container/:containerName", apis.ExecContainer)

}
