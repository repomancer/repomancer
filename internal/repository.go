package internal

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

type Repository struct {
	Host              string
	Organization      string
	Name              string
	BaseDir           string `json:"-"` // Calculated on load, not saved with configuration
	LogFile           string `json:"-"` // Calculated on load, not saved with configuration
	jobs              []*Job
	Status            string
	Selected          bool
	PullRequest       *PullRequest
	LastCommandResult error `json:"-"`
	mu                sync.Mutex
	RepositoryStatus  RepositoryStatus
}

func (r *Repository) GetUrl() *url.URL {
	repoUrl, _ := url.Parse(fmt.Sprintf("https://%s/%s/%s", r.Host, r.Organization, r.Name))
	return repoUrl
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
