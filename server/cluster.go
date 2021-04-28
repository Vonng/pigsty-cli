package server

import (
	"bufio"
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
)

// GetClusterHandler mount on GET /api/:cluster
func GetClusterHandler(c *gin.Context) {
	cluster := c.Param("cluster")
	cls := PS.Executor.Config.GetCluster(cluster)
	if cls == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("cluster %s not found", cluster),
		})
		return
	}

	c.JSON(http.StatusOK, cls)
	logrus.Infof("GET /api/v1/%s", cluster)
}

// PostClusterHandler will trigger cluster creation
func InitClusterHandler(c *gin.Context) {
	// get params and logo
	cluster := c.Param("cluster")
	force := c.Query("force")
	tags := c.QueryArray("tags")
	logrus.Infof("init cluster %s  force=%s tags=%s", cluster, force, tags)

	// create job with stdout pipe
	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()
	br := bufio.NewReaderSize(pr, 64)
	job := PS.Executor.NewJob(
		exec.WithPlaybook("pgsql.yml"),
		exec.WithName("pgsql init"),
		exec.WithLimit(cluster),
		exec.WithTags(tags...),
		exec.WithStdout(pw),
	)
	if force != "" && force != "false" {
		job.Opts.ExtraVars["pg_exists_action"] = "clean"
	}

	// create async reader generate sse event
	var wg sync.WaitGroup
	wg.Add(1)

	chanStream := make(chan string, 1)
	go func() {
		defer close(chanStream)
		for {
			s, err := br.ReadString('\n')
			if err != nil {
				break
			}
			chanStream <- s
			logrus.Infof(s)
			//c.SSEvent(cluster, s)
		}
		logrus.Infof("sse event sent")
		wg.Done()
	}()

	go func() {
		err := job.Run(c)
		pw.Close()
		if err != nil {
			logrus.Errorf("fail to run tasks %s", err)
		}
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			c.SSEvent(cluster, msg)
			return true
		}
		return false
	})

	//fmt.Println("run tasks, %v", job)

	//if err != nil {
	//	c.JSON(http.StatusConflict, gin.H{
	//		"message": err.Error(),
	//	})
	//} else {
	//	wg.Wait() // wait all message are flushed
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "ok",
	//	})
	//	return
	//}
}

// PostClusterHandler will trigger cluster creation
func RemoveClusterHandler(c *gin.Context) {
	// get params and logo
	cluster := c.Param("cluster")
	tags := c.QueryArray("tags")
	logrus.Infof("init cluster %s tags=%s", cluster, tags)

	// create job with stdout pipe
	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()
	br := bufio.NewReader(pr)
	job := PS.Executor.NewJob(
		exec.WithPlaybook("pgsql-remove.yml"),
		exec.WithName("pgsql remove"),
		exec.WithLimit(cluster),
		exec.WithTags(tags...),
		exec.WithStdout(pw),
	)

	// create async reader generate sse event
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			s, err := br.ReadString('\n')
			if err != nil {
				break
			}
			logrus.Infof(s)
			c.SSEvent(cluster, s)
		}
		logrus.Infof("sse event sent")
		wg.Done()
	}()

	//fmt.Println("run tasks, %v", job)
	err := job.Run(c)
	pw.Close()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	} else {
		wg.Wait() // wait all message are flushed
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		return
	}
}
