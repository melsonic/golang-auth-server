package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"message": "hello",
	})
}
