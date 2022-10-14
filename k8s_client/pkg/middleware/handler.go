package middleware

import "github.com/gin-gonic/gin"

// 跨域的头部
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		/*
			header("Access-Control-Allow-Origin:http://localhost:8080");
			header('Access-Control-Allow-Headers:token');//自定义请求头
			header("Access-Control-Allow-Credentials:true");
			header("Access-Control-Allow-Method:get,post");

		*/
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Method", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Headers", "authorization,origin,content-type,accept")

		c.AbortWithStatus(200)
	}

}
