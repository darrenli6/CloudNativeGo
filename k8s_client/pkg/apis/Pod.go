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
		return
	}
	c.JSON(http.StatusOK, pods)
}

func ExecContainer(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Param("containerName")
	method := c.DefaultQuery("method", "sh")
	err := service.WebSSH(namespaceName, podName, containerName, method, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

}
