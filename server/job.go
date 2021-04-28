package server

import (
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strings"
)

// GetConfigHandler will serve config file
func GetJobHandler(c *gin.Context) {
	if job := PS.GetJob(); job != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "job running",
			"data":    PS.Job,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "no job running",
			"data":    nil,
		})
	}
}

// ListJobHandler will iter job directory and return job list
func ListJobHandler(c *gin.Context) {
	if logInfo, err := PS.ListLogDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not list jobs",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "ok",
			"data":    logInfo,
		})
	}
}

// PostJobHandler will create new job
func PostJobHandler(c *gin.Context) {
	// arg parsing
	playbook := c.Query("playbook")
	cluster := c.Query("cluster")
	tags := c.QueryArray("tags")
	if !strings.HasSuffix(playbook, ".yml") {
		playbook += ".yml"
	}
	logrus.Infof("post job handler called: playbook=%s cluster=%s tags=%s", playbook, cluster, tags)

	// build new job
	job := PS.Executor.NewJob(
		exec.WithPlaybook(playbook),
		exec.WithName(fmt.Sprintf("%s-%s", playbook, cluster)),
		exec.WithLimit(cluster),
		exec.WithTags(tags...),
	)
	//logFilenName := fmt.Sprintf(`%s-%s@%s.log`)
	job.LogPath = filepath.Join(PS.LogDir(), job.ID)

	if j, err := PS.RunJob(job); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "job exists",
			"data":    j,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "job created",
			"data":    j,
		})
	}
	return
}

func DelJobHandler(c *gin.Context) {
	if job := PS.DelJob(); job != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "job deleted",
			"data":    job,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "job not found",
			"data":    nil,
		})
	}
}
