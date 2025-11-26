package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if c.Request.Header.Get("origin") != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": gin.H{},
	})
}

func main() {
	r := gin.Default()
	r.GET("/no_cors", Index)
	r.POST("/no_cors", Index)
	r.GET("/cors", Cors(), Index)
	r.POST("/cors", Cors(), Index)
	r.Run(":8080")
}
