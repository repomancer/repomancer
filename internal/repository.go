package internal

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
)

type Repository struct {
	Host              string
	Organization      string
	Name              string
	BaseDir           string    `json:"-"` // Calculated on load, not saved with configuration
	LogFile           string    `json:"-"` // Calculated on load, not saved with configuration
	Jobs              *JobQueue `json:"-"` // Created on load, not saved with configuration
	Selected          bool
	PullRequest       *PullRequest
	LastCommandResult error `json:"-"`
	RepositoryStatus  RepositoryStatus
	OnUpdated         func(repo *Repository) `json:"-"`
	JobsRunning       bool                   `json:"-"`
	jobMutex          sync.Mutex
}

func (r *Repository) changed() {
	if r.OnUpdated != nil {
		r.OnUpdated(r)
	} else {
		log.Printf("Repository %s has no OnUpdated", r.Name)
	}
}

func (r *Repository) AddJob(job *Job) {
	r.Jobs.Add(job)
	r.changed()
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

func (r *Repository) RunJobs() {
	if !r.jobMutex.TryLock() {
		// If jobs are already running in this repository, return. Any jobs that have been added will be picked up
		// on the same goroutine that's already running
		return
	}
	defer r.jobMutex.Unlock()
	r.JobsRunning = true
	for {
		job := r.Jobs.Pop()
		r.changed()
		if job == nil {
			break
		}
		job.Run()
	}
	r.JobsRunning = false
	r.changed()
}

func (r *Repository) JobStatus() string {
	var jobsString string
	if r.Jobs.Len() > 1 {
		jobsString = fmt.Sprintf("%d jobs pending", r.Jobs.Len())
	} else if r.Jobs.Len() == 1 {
		jobsString = "1 job pending"
	}
	if r.JobsRunning {
		if jobsString != "" {
			return fmt.Sprintf("Running (%s)", jobsString)
		} else {
			return "Running"
		}
	} else {
		return jobsString
	}
}

type RepositoryStatus struct {
	Cloned             bool
	BranchCreated      bool
	PullRequestCreated bool
	PullRequestClosed  bool
}
