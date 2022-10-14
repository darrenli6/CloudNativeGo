package apis

import (
	"net/http"

	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/service"
	"github.com/gin-gonic/gin"
)

func GetPods(c *gin.Context) {
	namespaceName := c.Param("namespaceName")

	pods, err := service.GetPod(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, pods)
}
