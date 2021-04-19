package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetClusterHandler mount on GET /api/:cluster
func GetClusterHandler(c *gin.Context) {
	cluster := c.Param("cluster")
	logrus.Infof("GET /api/pgsql/%s", cluster)
}
