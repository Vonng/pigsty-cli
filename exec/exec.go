package exec

import (
	"context"
	"github.com/Vonng/pigsty-cli/conf"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func (e *Executor) Reload(){}

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
	job.Status = JOB_READY
	e.Jobs[job.ID] = &job
	return &job
}

/**************************************************************\
*                          Job                                 *
\**************************************************************/
// Job is spawned by executor
type Job struct {
	ID       string // uuid v1
	Name     string // human readable job info
	Playbook string // playbook name
	Limit    string // limit execution targets
	Tags     []string
	Stdout   io.Writer // write output to this
	Stderr   io.Writer // write error to this
	Status   string    // ready | running | failed | success
	StartAt  time.Time
	DoneAt   time.Time
	CMD      *playbook.AnsiblePlaybookCmd     // ansible command
	Opts     *playbook.AnsiblePlaybookOptions // playbook options
	Exec     *Executor                        // Executor
}

// JobOpts will configure job
type JobOpts func(*Job)

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
	j.StartAt = time.Now()
	j.Status = JOB_RUNNING
	logrus.Infof(j.CMD.String())
	err := j.CMD.Run(ctx)
	j.DoneAt = time.Now()
	if err != nil {
		j.Status = JOB_FAILED
	} else {
		j.Status = JOB_SUCCESS
	}
	return err
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
