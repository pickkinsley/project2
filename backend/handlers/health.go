package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"version":    "1.0.0",
		"deployment": "automated",
	})
}
