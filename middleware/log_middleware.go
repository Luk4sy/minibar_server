package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

func LogMiddleware(c *gin.Context) {
	// 请求中间件
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())

	}
	fmt.Println("body:", string(byteData))

	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	c.Next()
	// 响应中间件
}
