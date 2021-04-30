package server

import (
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path"
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
		c.JSON(http.StatusOK, gin.H{
			"message": "job not found",
			"data":    nil,
		})
	}
}

// ListJobHandler will serve config file
func ListJobHandler(c *gin.Context) {
	if jobInfo, err := PS.LisJobDir(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "can not list jobs",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    jobInfo,
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
	job.LogPath = PS.LogPath(job.ID)
	logrus.Infof("new job created, log: %s", job.LogPath)

	// save job info to datadir/job/
	if err := PS.SaveJob(job); err != nil {
		logrus.Errorf("fail to save job to %s", PS.JobPath(job.ID))
	}

	if j, err := PS.RunJob(job); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "job exists",
			"data":    nil,
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
		c.JSON(http.StatusOK, gin.H{
			"message": "job not found",
			"data":    nil,
		})
	}
}

// ListLogHandler will iter log directory and return log(job) list
func ListLogHandler(c *gin.Context) {
	if logInfo, err := PS.ListLogDir(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "can not list logs",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    logInfo,
		})
	}
}

// Get log by job id
func GetLogHandler(c *gin.Context) {
	jobID := c.Param("jobid")
	logPath := PS.LogPath(jobID)
	b, _ := ioutil.ReadFile(logPath)
	c.String(http.StatusOK, string(b))
	//c.File(logPath)
}

// GetLatestLogHandler
func GetLatestLogHandler(c *gin.Context) {
	logs, err := PS.ListLogDir()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "fail to list jobs",
			"data":    nil,
		})
		return
	}

	var maxInd int
	var maxMtime int64
	for i, log := range logs {
		if log.Mtime >= maxMtime {
			maxInd = i
		}
	}

	latestLogName := logs[maxInd].Name
	jobPath := path.Join(PS.DataDir, latestLogName)
	logrus.Infof("server latest log %s of %v", jobPath, logs)
	c.File(jobPath)
}
