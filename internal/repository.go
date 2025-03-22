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
	BaseDir           string    `json:"-"` // Calculated on load, not saved with configuration
	LogFile           string    `json:"-"` // Calculated on load, not saved with configuration
	Jobs              *JobQueue `json:"-"` // Created on load, not saved with configuration
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
