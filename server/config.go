package server

import (
	"github.com/Vonng/pigsty-cli/conf"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetConfigHandler will serve config file
func GetConfigHandler(c *gin.Context) {
	c.File(PS.ConfigPath)
}

// PostConfigHandler will update default configuration file with posted content
// TODO: convenient but dangerous!!!
func PostConfigHandler(c *gin.Context) {
	d, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err = conf.ParseConfig(d)
	if err != nil {
		// invalid config
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := conf.OverwriteConfig(d, PS.ConfigPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := PS.Reload(PS.ConfigPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

//func OnConfigUpdate(data []byte) error {
//	_, err := conf.ParseConfig(data)
//	if err != nil {
//		// invalid config
//		c.JSON(http.StatusBadRequest, gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//	if err := conf.OverwriteConfig(d, PS.ConfigPath); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//
//	if err := PS.Reload(PS.ConfigPath); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "ok",
//	})
//}
