package main

import (
	"fmt"
	"github.com/gin-contrib/logger"
	"ibunny/tracer"
	"ibunny/tracker"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.DebugMode)
	// Logging to a file
	//f, err := os.Create("./logs/runtime.log")
	//if err != nil {
	//	panic("Cannot access log file.")
	//}
	//
	//gin.DefaultWriter = io.MultiWriter(f)
}

func main() {
	r := gin.Default()
	// Middlewares
	r.Use(logger.SetLogger())

	// Trusted proxies
	r.ForwardedByClientIP = true
	errSetProxies := r.SetTrustedProxies([]string{"127.0.0.1"})
	if errSetProxies != nil {
		return
	}

	// Homepage
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"page": "home",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		dt := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"now":     dt.String(),
		})
	})

	// Debug
	r.GET("/debug", func(c *gin.Context) {
		dt := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"code":      "SUCCESS",
			"now":       dt.String(),
			"ip":        c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
			"referer":   c.Request.Referer(),
		})
	})

	// Redirect by code
	r.GET("/:code", tracker.ShortCodeHandler)

	// Redirect by structural
	r.GET("/go/:id", tracker.GoHandler)

	// Trace URL
	r.POST("/trace", tracer.TraceHandler)

	errRun := r.Run("127.0.0.1:8000")
	if errRun != nil {
		fmt.Println("Error when start server.")
		return
	}
}
