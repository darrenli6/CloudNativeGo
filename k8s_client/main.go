package main

import (
	"fmt"

	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/config"
	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/router"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func main() {

	engine := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.InitRouter(engine)
	err := engine.Run(fmt.Sprintf("%s:%d", config.GetString(config.ServerHost), config.GetInt(config.ServerPort)))
	if err != nil {
		klog.Fatal(err)
		return
	}

}
