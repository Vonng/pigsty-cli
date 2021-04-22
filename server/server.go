package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
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

var PS *PigstyServer

type PigstyServer struct {
	ListenAddr string
	ConfigPath string
	PublicDir  string
	HomeDir    string
	Server     *http.Server
	Executor   *exec.Executor
}

func (ps *PigstyServer) ServerDir() string {
	return filepath.Join(ps.HomeDir, ".pigsty")
}

func (ps *PigstyServer) ResourceDir() string {
	// return filepath.Join(ps.HomeDir, ".pigsty", "public")
	return "./public"
}

// NewPigstyServer will create new server
func NewPigstyServer(configPath string, publicDir string, listenAddr string) *PigstyServer {
	var ps PigstyServer
	ps.ListenAddr = listenAddr
	ps.ConfigPath = configPath
	ps.PublicDir = publicDir
	if ps.Executor = exec.NewExecutor(configPath); ps.Executor == nil {
		return nil
	}
	ps.HomeDir = ps.Executor.WorkDir

	// build router
	r := gin.Default()

	/******************************************
	 * config interface
	 ******************************************/
	// get config, put config, validate config
	r.GET("/api/v1/config", GetConfigHandler)
	r.POST("/api/v1/config", PostConfigHandler)
	r.GET("/api/v1/pgsql/:cluster/info", GetClusterHandler)
	r.GET("/api/v1/pgsql/:cluster/init", InitClusterHandler)
	r.GET("/api/v1/pgsql/:cluster/remove", RemoveClusterHandler)
	//r.GET("/api/v1/pgsql/:cluster/init", InitClusterHandler)
	//r.GET("/api/v1/pgsql/:cluster/remove", RemoveClusterHandler)

	/******************************************
	 * static resource
	 ******************************************/
	//logrus.Infof("serve resource from %s", ps.ResourceDir())
	r.Use(static.Serve("/", static.LocalFile(publicDir, true)))

	srv := &http.Server{
		Addr:    ps.ListenAddr,
		Handler: r,
	}

	// return PigstyServer
	ps.Server = srv
	return &ps
}

// Reload will create a new Executor according to config
func (ps *PigstyServer) Reload(configPath string) error {
	executor := exec.NewExecutor(configPath)
	if executor == nil {
		return fmt.Errorf("reload failed: invalid config")
	}
	ps.Executor = executor
	return nil
}

func (ps *PigstyServer) Run() {
	go func() {
		if err := ps.Server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ps.Server.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}
	logrus.Println("Server exiting")
}

func InitDefaultServer(configPath string, publicDir string, listenAddr string) {
	PS = NewPigstyServer(configPath, publicDir, listenAddr)
	PS.Run()
}
