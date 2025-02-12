package internal

import (
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"strings"
	"sync"
	"time"
)

type Repository struct {
	Host              string
	Organization      string
	Name              string
	BaseDir           string `json:"-"` // Calculated on load, not saved with configuration
	jobs              []*Job
	Status            string
	Selected          bool
	PullRequest       *PullRequest
	LastCommandResult error
	mu                sync.Mutex
	RepositoryStatus  RepositoryStatus
	logBinding        binding.String
	log               []string
}

func (r *Repository) GetLogBinding() binding.String {
	if r.logBinding == nil {
		r.logBinding = binding.NewString()
		_ = r.logBinding.Set(strings.Join(r.log, "\n"))
	}
	return r.logBinding
}

func (r *Repository) Log(message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	withTimestamp := fmt.Sprintf("[%s]: %s", time.Now().Format("2006-01-02 15:04:05"), message)
	r.log = append(r.log, withTimestamp)
	if r.logBinding != nil {
		_ = r.logBinding.Set(strings.Join(r.log, "\n"))
	}
}

type PullRequest struct {
	Number                 int
	Url                    string
	Status                 string
	LastChecked            time.Time
	StatusCheckRollupState string
}

func (r *Repository) Title() string {
	return fmt.Sprintf("%s/%s/%s", r.Host, r.Organization, r.Name)
}

type RepositoryStatus struct {
	Cloned             bool
	BranchCreated      bool
	PullRequestCreated bool
	PullRequestClosed  bool
}

func (r *Repository) AddJob(job *Job) {
	r.mu.Lock()
	r.jobs = append(r.jobs, job)
	r.mu.Unlock()
}

func (r *Repository) GetJob(i int) *Job {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.jobs[i]
}

func (r *Repository) JobCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.jobs)
}

func (r *Repository) QueuedJobs() int {
	cnt := 0
	for i := 0; i < len(r.jobs); i++ {
		if !r.jobs[i].Finished {
			cnt++
		}
	}
	return cnt
}
