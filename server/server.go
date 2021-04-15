package main

import (
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

/*************************************************************\
*                          Routers                            *
GET    /api/config                   get global config file
POST   /api/config                   update global config file
POST   /api/:cluster                 init new pgsql cluster <cluster>
POST   /api/:cluster/:instance       init new pgsql <instance> under <cluster>
DELETE /api/:cluster                 destroy pgsql cluster <cluster>                   {job_id}
DELETE /api/:instance                destroy pgsql instance <instance>                 {job_id}
GET    /api/infra                    get global infra info
GET    /api/pgsql                    get all cluster info
GET    /api/:cluster                 get cluster info of <cluster>
GET    /api/:cluster/jobs            # get jobs of cluster                            # SSE
GET    /api/:cluster/:job            # get job info                                   # SSE
POST   /api/infra/targets            update filesd config
POST   /api/infra/haproxy            update haproxy


GET   /api/                          get API list

GET  /api/:cluster/instances         get instances info of <cluster>
GET  /api/:cluster/services          get services info of <cluster>
GET  /api/:cluster/users             get users info of <cluster>
GET  /api/:cluster/databases         get databases info of <cluster>
GET  /api/:cluster/primary           get primary instance info of <cluster>
GET  /api/:cluster/:instance         get instance info of <cluster>
GET  /api/:cluster/users             get databases info of <cluster>
GET  /api/:cluster/databases         get databases info of <cluster>
POST /api/:cluster                   create new pgsql cluster <cluster>
POST /api/:cluster/:instance         create new pgsql <instance> under <cluster>
POST /api/:cluster/:user             create new biz user under <cluster>
POST /api/:cluster/:database         create new biz database under <cluster>
POST /api/:cluster/:hba              create new hba rules
POST /api/:cluster/:job              create new hba rules
\*************************************************************/

var EX = exec.NewExecutor("/Users/vonng/pigsty/pigsty.yml")

func NewEngine() *gin.Engine {
	//gin.DefaultWriter = colorable.NewColorableStderr()
	r := gin.Default()

	/******************************************
	* config interface
	 ******************************************/
	// get config, put config, validate config
	r.GET("/api/config", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, EX.Config)
	})
	r.POST("/api/config", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, EX.Config)
	})
	r.POST("/api/config/validate", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, EX.Config)
	})

	/******************************************
	* get generate information
	 ******************************************/
	// get clusters list info
	r.GET("/api/pgsql", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, EX.Config.Clusters)
	})

	// get infrastructure info
	r.GET("/api/infra", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, EX.Config.Clusters)
	})

	/******************************************
	* cluster management
	 ******************************************/
	// Single cluster:  get=info | post=create | delete=remove
	// get: cluster info
	r.GET("/api/pgsql/:cluster", func(c *gin.Context) {
		logrus.Infof("select cluster %s", c.Param("cluster"))
		c.PureJSON(http.StatusOK, EX.Config.GetCluster(c.Param("cluster")))
	})
	r.POST("/api/pgsql/:cluster", func(c *gin.Context) {
		logrus.Infof("create cluster %s", c.Param("cluster"))
		c.PureJSON(http.StatusOK, EX.Config.Clusters)
	})
	r.DELETE("/api/pgsql/:cluster", func(c *gin.Context) {
		logrus.Infof("remove cluster %s", c.Param("cluster"))
		c.PureJSON(http.StatusOK, EX.Config.Clusters)
	})

	/******************************************
	* instance management
	 ******************************************/
	r.GET("/api/pgsql/:cluster/seq/:seq", func(c *gin.Context) {
		cluster, seq := c.Param("cluster"), c.Param("seq")
		instance := fmt.Sprintf("%s-%s", cluster, seq)
		logrus.Infof("create instance %s", instance)
		c.PureJSON(http.StatusOK, EX.Config.GetInstance(instance))
	})
	r.POST("/api/:cluster/:seq", func(c *gin.Context) {
		logrus.Infof("create instance %s-%s", c.Param("cluster"), c.Param("seq"))
		c.PureJSON(http.StatusOK, EX.Config.Clusters)
	})
	r.DELETE("/api/:cluster/:seq", func(c *gin.Context) {
		logrus.Infof("remove instance %s-%s", c.Param("cluster"), c.Param("seq"))
		c.PureJSON(http.StatusOK, EX.Config.Clusters)
	})

	r.GET("/stream", func(c *gin.Context) {
		chanStream := make(chan int, 10)
		go func() {
			defer close(chanStream)
			for i := 0; i < 500; i++ {
				chanStream <- i
				time.Sleep(time.Second * 1)
			}
		}()
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-chanStream; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	return r
}

func main() {
	r := NewEngine()
	r.Run()
}
