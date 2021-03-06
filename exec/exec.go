package exec

import (
	"context"
	"encoding/json"
	"github.com/Vonng/pigsty-cli/conf"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/**************************************************************\
*                          Const                               *
\**************************************************************/
const (
	JOB_READY   = "ready"
	JOB_RUNNING = "running"
	JOB_FAILED  = "failed"
	JOB_SUCCESS = "success"
)

/**************************************************************\
*                        Executor                              *
\**************************************************************/
// Executor hold config path and content
type Executor struct {
	WorkDir   string
	Inventory string
	Config    *conf.Config
	Jobs      map[string]*Job
	Lock      *sync.Mutex
}

// NewExecutor will create ansible playbook executor based on config path
func NewExecutor(path string) *Executor {
	configPath, err := filepath.Abs(path)
	if err != nil {
		logrus.Fatalf("invalid config path %s, %s", path, err)
		return nil
	}

	pigstyFile, pigstyDir := filepath.Base(configPath), filepath.Dir(configPath)
	fi, err := os.Stat(configPath)
	if err != nil {
		logrus.Fatalf("invalid inventory path %s, %s", configPath, err)
		return nil
	}
	if fi.IsDir() { // if dir is given, assume pigsty home dir
		pigstyFile, pigstyDir = "pigsty.yml", configPath
		configPath = filepath.Join(configPath, `pigsty.yml`)
		if fi, err = os.Stat(configPath); err != nil {
			logrus.Fatalf("could not find pigsty.yml in %s, %s", pigstyDir, err)
			return nil
		}
	}

	cfg, err := conf.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("fail to parse config %s, %s", path, err)
		return nil
	}

	logrus.Debugf("load config from %s", path)
	return &Executor{
		WorkDir:   pigstyDir,
		Inventory: pigstyFile,
		Config:    cfg,
		Jobs:      make(map[string]*Job),
	}
}

func (e *Executor) Reload() {}

// Static return static resource of this executor (.pigsty/public by default)
func (e *Executor) StaticDir() string {
	return filepath.Join(e.WorkDir, ".pigsty", "public")
}

// LogDir return log directory of this executor (.pigsty/job by default)
func (e *Executor) LogDir() string {
	return filepath.Join(e.WorkDir, ".pigsty", "log")
}

// NewJob will spawn new job and modify it with JobOpts
func (e *Executor) NewJob(options ...JobOpts) *Job {
	var job Job
	job.Opts = &playbook.AnsiblePlaybookOptions{ExtraVars: map[string]interface{}{}}
	job.Exec = e
	id, err := uuid.NewUUID()
	if err != nil {
		id = uuid.New()
	}
	job.ID = id.String()
	for _, opt := range options {
		opt(&job)
	}

	// overwrite important options
	if job.Limit != "" && job.Opts.Limit == "" {
		job.Opts.Limit = job.Limit
	}
	if job.Tags != nil && len(job.Tags) > 0 && job.Opts.Tags == "" {
		job.Opts.Tags = strings.Join(job.Tags, ",")
	}
	execOpts := []execute.ExecuteOptions{execute.WithCmdRunDir(e.WorkDir)}
	if job.Stdout != nil {
		execOpts = append(execOpts, execute.WithWrite(job.Stdout))
	}
	if job.Stderr != nil {
		execOpts = append(execOpts, execute.WithWriteError(job.Stderr))
	}
	job.CMD = &playbook.AnsiblePlaybookCmd{
		Playbooks: []string{job.Playbook},
		Options:   job.Opts,
		Exec:      execute.NewDefaultExecute(execOpts...),
	}
	job.Command = job.CMD.String()
	job.StartAt = time.Now()
	job.Status = JOB_READY
	e.Jobs[job.ID] = &job
	return &job
}

/**************************************************************\
*                          Job                                 *
\**************************************************************/
// Job is spawned by executor
type Job struct {
	ID       string                           `json:"id"`       // uuid v1
	Name     string                           `json:"name"`     // human readable job info
	Playbook string                           `json:"playbook"` // playbook name
	Limit    string                           `json:"limit"`    // limit execution targets
	Tags     []string                         `json:"tags"`     // execution tags
	LogPath  string                           `json:"log_path"` // write playbook log to ANSIBLE_LOG_PATH
	Status   string                           `json:"status"`   // ready | running | failed | success
	StartAt  time.Time                        `json:"start_at"` // job start at
	DoneAt   time.Time                        `json:"done_at"`  // job done at
	Command  string                           `json:"command"`  // job raw shell command
	CMD      *playbook.AnsiblePlaybookCmd     `json:"-"`        // ansible command
	Opts     *playbook.AnsiblePlaybookOptions `json:"-"`        // playbook options
	Exec     *Executor                        `json:"-"`        // Executor

	Stdout io.Writer          `json:"-"` // write output to this
	Stderr io.Writer          `json:"-"` // write error to this
	ctx    context.Context    // job context
	cancel context.CancelFunc // job cancel func
}

// JobOpts will configure job
type JobOpts func(*Job)

func (j *Job) JSON() string {
	b, _ := json.Marshal(j)
	return string(b)
}

// WithStdout will set stdout
func WithStdout(w io.Writer) JobOpts {
	return func(j *Job) {
		j.Stdout = w
	}
}

// WithStderr will set stderr output
func WithStderr(w io.Writer) JobOpts {
	return func(j *Job) {
		j.Stderr = w
	}
}

// WithTags will set job name
func WithName(name string) JobOpts {
	return func(j *Job) {
		j.Name = name
	}
}

// WithTags will set playbook names
func WithPlaybook(playbook string) JobOpts {
	return func(j *Job) {
		j.Playbook = playbook
	}
}

// WithTags will set tasks tags to playbook
func WithTags(tags ...string) JobOpts {
	return func(j *Job) {
		j.Tags = tags
	}
}

// WithLimit will set limit string to playbook
func WithLimit(limit string) JobOpts {
	return func(j *Job) {
		j.Limit = limit
	}
}

func WithLogPath(logPath string) JobOpts {
	return func(j *Job) {
		j.LogPath = logPath
	}
}

// WithAnsibleOpts will overwrite entire options, use with caution!
func WithAnsibleOpts(opts *playbook.AnsiblePlaybookOptions) JobOpts {
	return func(j *Job) {
		j.Opts = opts
	}
}

// WithExtraVars will add extra k-v into playbook's extra vars
func WithExtraVars(key string, value interface{}) JobOpts {
	return func(j *Job) {
		if j.Opts == nil {
			j.Opts = &playbook.AnsiblePlaybookOptions{ExtraVars: map[string]interface{}{}}
		} else {
			if j.Opts.ExtraVars == nil {
				j.Opts.ExtraVars = map[string]interface{}{}
			}
		}
		j.Opts.ExtraVars[key] = value
	}
}

// CreateLog will create job log under executor's log directory
func (j *Job) CreateLog(name string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(j.Exec.LogDir(), name))
}

// Run will run given command under context
func (j *Job) Run(ctx context.Context) error {
	setupOsEnv()
	// create filelog
	j.Command = j.CMD.String()
	if j.LogPath != "" {
		_ = os.Setenv("ANSIBLE_LOG_PATH", j.LogPath)
		f, err := os.Create(j.LogPath)
		f.Close()
		if err != nil {
			j.Status = JOB_FAILED
			return err
		}
	}
	j.StartAt = time.Now()
	j.ctx, j.cancel = context.WithCancel(ctx)
	j.Status = JOB_RUNNING
	logrus.Infof(j.CMD.String())
	err := j.CMD.Run(j.ctx)
	j.DoneAt = time.Now()
	if err != nil {
		logrus.Errorf("job failed: %s", err)
		j.Status = JOB_FAILED
	} else {
		j.Status = JOB_SUCCESS
	}
	return err
}

func (j *Job) Cancel() {
	if j.cancel != nil {
		j.cancel()
	}
	return
}

func (j *Job) MutexRun(ctx context.Context) error {
	j.Exec.Lock.Lock()
	defer j.Exec.Lock.Unlock()
	return j.Run(ctx)
}

// AsyncRun will run tasks on background, a job id is returned
func (j *Job) AsyncRun() string {
	go func() {
		if j.Stdout == nil {
			f, err := j.CreateLog(j.ID + ".log")
			if err != nil {
				logrus.Errorf("fail to create log of job %v", j)
				return
			}
			defer f.Close()
			j.Stdout = f
		}

		if err := j.Run(context.TODO()); err != nil {
			logrus.Errorf("job failed: %v", j)
		}
	}()
	return j.ID
}

// setupOsEnv will setup ansible environment
func setupOsEnv() {
	_ = os.Setenv("ANSIBLE_TRANSFORM_INVALID_GROUP_CHARS", "ignore")
	_ = os.Setenv("ANSIBLE_DEPRECATION_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_SYSTEM_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_ACTION_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_COMMAND_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_DEVEL_WARNING", "False")
	_ = os.Setenv("ANSIBLE_DISPLAY_ARGS_TO_STDOUT", "False")
}
