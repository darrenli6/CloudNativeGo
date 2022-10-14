package apis

import (
	"net/http"

	"github.com/darrenli6/CloudNativeGo/k8s_client/pkg/service"
	"github.com/gin-gonic/gin"
)

func GetNamespace(c *gin.Context) {

	namespace, err := service.GetNamespace()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, namespace)
}
