package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// PS is the default PigstyServer
var PS *PigstyServer

type PigstyServer struct {
	ListenAddr string
	ConfigPath string
	PublicDir  string
	HomeDir    string
	Server     *http.Server
	Executor   *exec.Executor
	Job        *exec.Job // shitty-implementation: only one job allow one time
	lock       sync.Mutex
	jobLock    sync.RWMutex
	cancel     context.CancelFunc
}

// NewPigstyServer will create new server
func NewPigstyServer(configPath string, publicDir string, listenAddr string) *PigstyServer {
	var ps PigstyServer
	ps.ListenAddr = listenAddr
	ps.ConfigPath = configPath
	ps.PublicDir = publicDir

	if fi, err := os.Stat(ps.PublicDir); err != nil && os.IsNotExist(err) && !fi.IsDir() {
		logrus.Errorf("public dir %s not exists", ps.PublicDir)
		return nil
	}
	// make sure log dir exists
	_ = os.Mkdir(ps.LogDir(), 0755)

	if ps.Executor = exec.NewExecutor(configPath); ps.Executor == nil {
		return nil
	}
	ps.HomeDir = ps.Executor.WorkDir
	ps.Server = &http.Server{
		Addr:    ps.ListenAddr,
		Handler: DefaultRouter(publicDir),
	}
	return &ps
}

func (ps *PigstyServer) LogDir() string {
	return filepath.Join(ps.PublicDir, "log")
}

type LogInfo struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Mtime int64  `json:"mtime"`
}

// ListLogdir will return
func (ps *PigstyServer) ListLogDir() ([]LogInfo, error) {
	logs, err := ioutil.ReadDir(ps.LogDir())
	if err != nil {
		return nil, err
	}
	var list []LogInfo
	for _, log := range logs {
		if !log.IsDir() {
			list = append(list, LogInfo{log.Name(), log.Size(), log.ModTime().Unix()})
		}
	}
	return list, nil
}

// Reload will create a new Executor according to config
func (ps *PigstyServer) Reload(configPath string) error {
	if ps.Job != nil {
		return fmt.Errorf("executor can not be reloaed while running job")
	}
	// acquire lock
	ps.lock.Lock()
	defer ps.lock.Unlock()
	// TODO: can not reload while running jobs
	executor := exec.NewExecutor(configPath)
	if executor == nil {
		return fmt.Errorf("reload failed: invalid config")
	}
	ps.Executor = executor
	ps.HomeDir = configPath
	return nil
}

// RunJob will run job on background, error if running job already exists
func (ps *PigstyServer) RunJob(job *exec.Job) (*exec.Job, error) {
	ps.jobLock.Lock()
	defer ps.jobLock.Unlock()
	if ps.Job != nil && (ps.Job.Status == exec.JOB_RUNNING || ps.Job.Status == exec.JOB_READY) {
		return ps.Job, fmt.Errorf("another job is running: %s", ps.Job.ID)
	}
	ps.Job = job
	ctx, cancel := context.WithCancel(context.TODO())
	ps.cancel = cancel
	go ps.Job.Run(ctx)
	return job, nil
}

// CancelJob will cancel current job
func (ps *PigstyServer) DelJob() *exec.Job {
	ps.jobLock.Lock()
	defer ps.jobLock.Unlock()

	if ps.Job != nil && (ps.Job.Status == exec.JOB_RUNNING || ps.Job.Status == exec.JOB_READY) {
		ps.Job.Cancel()
		job := ps.Job
		ps.Job = nil
		return job
	}
	if ps.Job != nil {
		ps.Job = nil
	}
	return nil
}

// Get Job will return current running job, nil if not exists
func (ps *PigstyServer) GetJob() *exec.Job {
	ps.jobLock.RLock()
	defer ps.jobLock.RUnlock()
	if ps.Job != nil {
		if ps.Job.Status == exec.JOB_RUNNING || ps.Job.Status == exec.JOB_READY {
			return ps.Job
		} else {
			ps.Job = nil // remove finished jobs
			return nil
		}
	}
	return nil
}

func DefaultRouter(publicDir string) *gin.Engine {
	r := gin.Default()
	// config
	r.GET("/api/v1/config/", GetConfigHandler)
	r.POST("/api/v1/config/", PostConfigHandler)

	// job (list get create del)
	r.GET("/api/v1/jobs", ListJobHandler)
	r.GET("/api/v1/job", GetJobHandler)
	r.POST("/api/v1/job", PostJobHandler)
	r.DELETE("/api/v1/job", DelJobHandler)
	r.Use(static.Serve("/", static.LocalFile(publicDir, true)))
	return r
}

// run will launch server and listen
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

	// cancel running job
	if ps.Job != nil {
		ps.Job.Cancel()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ps.Server.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}
	logrus.Println("Server exiting")
}

// InitDefaultServer will init default pigsty singleton
func InitDefaultServer(configPath string, publicDir string, listenAddr string) {
	PS = NewPigstyServer(configPath, publicDir, listenAddr)
	if PS == nil {
		os.Exit(1)
	}
	PS.Run()
}
