package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(c *gin.Context) {
	start := time.Now()

	c.Next()

	requestMethod := c.Request.Method
	requestPath := c.Request.URL.Path
	requestAddress := c.Request.RemoteAddr
	requestOrigin := c.GetHeader("Origin")
	requestUserAgent := c.GetHeader("User-Agent")

	status := c.Writer.Status()
	latency := time.Since(start).Milliseconds()

	now := time.Now()
	formatted := now.Format("01/02/06 15:04:05")
	fmt.Println()
	fmt.Printf("	Request @ %s", formatted)
	fmt.Printf("		STATUS: %d @ %dms\n", status, latency)
	fmt.Printf("		METHOD: %s %s\n", requestMethod, requestPath)
	fmt.Printf("		Origin: %s\n", requestOrigin)
	fmt.Printf("		IP: %s\n", requestAddress)
	fmt.Printf("		User-Agent, %s", requestUserAgent)
	fmt.Println()
}
