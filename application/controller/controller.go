package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartGin(port string, r *gin.Engine) {
	SystemRoutes(r)
	AppRoutes(r)
	r.Run(port)
}

func SystemRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, fc!")
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
