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

const (
	JOB_READY   = "ready"
	JOB_RUNNING = "running"
	JOB_FAILED  = "failed"
	JOB_SUCCESS = "success"
)

// Jobs record all jobs
var Jobs map[string]*Job = make(map[string]*Job)

// Executor hold config path and content
type Executor struct {
	WorkDir   string
	Inventory string
	Config    *conf.Config
}

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

func WithStdout(w io.Writer) JobOpts {
	return func(j *Job) {
		j.Stdout = w
	}
}

func WithStderr(w io.Writer) JobOpts {
	return func(j *Job) {
		j.Stderr = w
	}
}

func WithName(name string) JobOpts {
	return func(j *Job) {
		j.Name = name
	}
}

func WithPlaybook(playbook string) JobOpts {
	return func(j *Job) {
		j.Playbook = playbook
	}
}

func WithTags(tags ...string) JobOpts {
	return func(j *Job) {
		j.Tags = tags
	}
}

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

// NewExecutor will create ansible playbook executor based on config path
func NewExecutor(path string) *Executor {
	configPath, err := filepath.Abs(path)
	if err != nil {
		logrus.Fatalf("invalid config path %s, %w", path, err)
		return nil
	}

	pigstyFile, pigstyDir := filepath.Base(configPath), filepath.Dir(configPath)
	fi, err := os.Stat(configPath)
	if err != nil {
		logrus.Fatalf("invalid inventory path %s, %w", configPath, err)
		return nil
	}
	if fi.IsDir() { // if dir is given, assume pigsty home dir
		pigstyFile, pigstyDir = "pigsty.yml", configPath
		configPath = filepath.Join(configPath, `pigsty.yml`)
		if fi, err = os.Stat(configPath); err != nil {
			logrus.Fatalf("could not find pigsty.yml in %s, %w", pigstyDir, err)
			return nil
		}
	}

	cfg, err := conf.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("fail to parse config %s, %w", path, err)
		return nil
	}
	logrus.Debugf("load config from %s", path)
	return &Executor{
		WorkDir:   pigstyDir,
		Inventory: pigstyFile,
		Config:    cfg,
	}
}

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
	Jobs[job.ID] = &job
	return &job
}

// Run will run given command under context
func (j *Job) Run(ctx context.Context) error {
	_ = os.Setenv("ANSIBLE_TRANSFORM_INVALID_GROUP_CHARS", "ignore")
	_ = os.Setenv("ANSIBLE_DEPRECATION_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_SYSTEM_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_ACTION_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_COMMAND_WARNINGS", "False")
	_ = os.Setenv("ANSIBLE_DEVEL_WARNING", "False")
	_ = os.Setenv("ANSIBLE_DISPLAY_ARGS_TO_STDOUT", "False")

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


// Async
func (j *Job) AsyncRun() string {
	return ""
}
