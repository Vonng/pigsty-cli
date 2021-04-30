package server

import (
	"context"
	"embed"
	"encoding/json"
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

/****************************************************
*  Embed Resources
/****************************************************/

//go:embed index.html
//go:embed static
//go:embed img
//go:embed *.ico
//go:embed *.png
//go:embed *.json
//go:embed *.txt
var Resource embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func EmbedFileSystem(fs http.FileSystem) *embedFileSystem {
	return &embedFileSystem{fs}
}

func (fs *embedFileSystem) Exists(prefix string, filepath string) bool {
	return true
}

/****************************************************
*  Pigsty Server
/****************************************************/
// PS is the default PigstyServer
var PS *PigstyServer

// PigstyServer holds required information
type PigstyServer struct {
	ListenAddr string
	ConfigPath string
	DataDir    string
	PublicDir  string
	HomeDir    string
	Server     *http.Server
	Executor   *exec.Executor
	Job        *exec.Job // shitty-implementation: only one job allow one time
	lock       sync.Mutex
	jobLock    sync.RWMutex
	cancel     context.CancelFunc
}

// ServerOpt will configure pigsty server
type ServerOpt func(server *PigstyServer)

// WithStdout will set stdout
func WithListenAddress(listenAddr string) ServerOpt {
	return func(ps *PigstyServer) {
		ps.ListenAddr = listenAddr
	}
}

func WithPublicDir(publicDir string) ServerOpt {
	return func(ps *PigstyServer) {
		ps.PublicDir = publicDir
	}
}

func WithDataDir(dataDir string) ServerOpt {
	return func(ps *PigstyServer) {
		ps.DataDir = dataDir
	}
}

func WithConfigPath(configPath string) ServerOpt {
	return func(ps *PigstyServer) {
		ps.ConfigPath = configPath
	}
}

// NewPigstyServer will create new server
func NewPigstyServer(opts ...ServerOpt) *PigstyServer {
	var ps PigstyServer
	for _, opt := range opts {
		opt(&ps)
	}
	// make sure data (log|job) dir exists
	if err := MakeDataDir(ps.DataDir); err != nil {
		logrus.Fatalf("fail to log dir %s %s", ps.DataDir, err.Error())
		return nil
	}
	if ps.Executor = exec.NewExecutor(ps.ConfigPath); ps.Executor == nil {
		return nil
	}
	ps.HomeDir = ps.Executor.WorkDir
	ps.Server = &http.Server{
		Addr:    ps.ListenAddr,
		Handler: ps.DefaultRouter(),
	}
	return &ps
}

/****************************************************
*  Log Info
/****************************************************/
// MakeDataDir will make sure data dir exists
func MakeDataDir(path string) error {
	logrus.Infof("check data dir %s", path)
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		// if path exists and is a file (exit)
		return fmt.Errorf("path exists and is regular file: %s %w", path, err)
	}
	_ = os.MkdirAll(path, 0755)
	_ = os.MkdirAll(filepath.Join(path, "log"), 0755)
	_ = os.MkdirAll(filepath.Join(path, "job"), 0755)
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return nil
	} else {
		return fmt.Errorf("fail to make data dir: %s", path)
	}
}

func (ps *PigstyServer) JobDir() string {
	return filepath.Join(ps.DataDir, "job")
}

func (ps *PigstyServer) LogDir() string {
	return filepath.Join(ps.DataDir, "log")
}

func (ps *PigstyServer) LogPath(name string) string {
	return filepath.Join(ps.DataDir, "log", name)
}

func (ps *PigstyServer) JobPath(name string) string {
	return filepath.Join(ps.DataDir, "job", name)
}

func (ps *PigstyServer) SaveJob(job *exec.Job) error {
	jobPath := ps.JobPath(job.ID)
	b, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(jobPath, b, 0644)
}

// LogInfo Hold log name, size, mtime info of job logs
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

func (ps *PigstyServer) LisJobDir() ([]exec.Job, error) {
	jobs, err := ioutil.ReadDir(ps.JobDir())
	if err != nil {
		return nil, err
	}
	var jobData []exec.Job
	for _, jobID := range jobs {
		job := ps.LoadJob(jobID.Name())
		if job != nil {
			jobData = append(jobData, *job)
		}
	}
	return jobData, nil
}

func (ps *PigstyServer) LoadJob(id string) *exec.Job {
	jobPath := ps.JobPath(id)
	b, err := ioutil.ReadFile(jobPath)
	if err != nil {
		return nil
	}
	var job exec.Job
	if err = json.Unmarshal(b, &job); err != nil {
		return nil
	}
	return &job
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

func (ps *PigstyServer) DefaultRouter() *gin.Engine {
	r := gin.Default()
	// config
	r.GET("/api/v1/config", GetConfigHandler)
	r.GET("/api/v1/config/", GetConfigHandler)
	r.POST("/api/v1/config/", PostConfigHandler)
	r.POST("/api/v1/config", PostConfigHandler)

	// job (get post del)
	r.GET("/api/v1/jobs", ListJobHandler)
	r.GET("/api/v1/job", GetJobHandler)
	r.POST("/api/v1/job", PostJobHandler)
	r.DELETE("/api/v1/job", DelJobHandler)

	// log (list latest get)
	r.GET("/api/v1/log/", ListLogHandler)
	r.GET("/api/v1/log/latest", GetLatestLogHandler)
	r.GET("/api/v1/log/:jobid", GetLogHandler)

	// use embed static resource or public dir if specified
	if ps.PublicDir == "" || ps.PublicDir == "embed" {
		logrus.Infof("use embed public resource")
		r.Use(static.Serve("/", EmbedFileSystem(http.FS(Resource))))
	} else {
		logrus.Infof("use public dir resource @ %s", ps.PublicDir)
		r.Use(static.Serve("/", static.LocalFile(ps.PublicDir, true)))
	}

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
func InitDefaultServer(listenAddr, configPath, dataDir, publicDir string) {
	logrus.Infof("pigsty server listen on %s , pigsty-config=%s  , dataDir=%s, publicDir=%s", listenAddr, configPath, dataDir, publicDir)
	PS = NewPigstyServer(
		WithListenAddress(listenAddr),
		WithDataDir(dataDir),
		WithPublicDir(publicDir),
		WithConfigPath(configPath),
	)
	if PS == nil {
		os.Exit(1)
	}
	os.Setenv("GIN_MODE", "release")
	gin.SetMode(gin.ReleaseMode)
	PS.Run()
}
