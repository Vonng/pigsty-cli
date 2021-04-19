package main

import (
	"github.com/gin-gonic/gin"
)

// GetConfigHandler mount on GET /api/:cluster
func GetConfigHandler(c *gin.Context) {
	c.File("/Users/vonng/pigsty/pigsty.yml")
}

func PutConfigHandler(c *gin.Context) {
	c.File("/Users/vonng/pigsty/pigsty.yml")
}

func PostConfigHandler(c *gin.Context) {
	c.File("/Users/vonng/pigsty/pigsty.yml")
}
