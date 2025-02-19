package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"time"
)

type RepositoryInfo struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	PushedAt string `json:"pushedAt"`
}

const PullRequestFilename = "PullRequest.md"

// NormalizeGitUrl tries to take any input URL and turn it to
// github-host.com/org/repo
// There are a lot of corner cases that don't work here, it is not an exhaustive list
// For example, it doesn't try to clean out invalid characters
func NormalizeGitUrl(url string) (string, error) {
	u := strings.TrimPrefix(url, "ssh://")
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "http://")
	u = strings.TrimPrefix(u, "git://")
	u = strings.TrimSuffix(u, ".git/")
	u = strings.TrimSuffix(u, ".git")
	s := strings.Split(u, "/")
	if len(s) < 3 || len(strings.TrimSpace(s[0])) == 0 || len(strings.TrimSpace(s[1])) == 0 || len(strings.TrimSpace(s[2])) == 0 {
		return "", fmt.Errorf("invalid git url: %s", url)
	}
	return u, nil
}

func Clone(r *Repository) (string, error) {
	cmd := fmt.Sprintf("gh repo clone %s/%s/%s . -- --depth=1", r.Host, r.Organization, r.Name)
	stdout, stderr, err := ShellOut(cmd, r.BaseDir)
	if err != nil {
		return "", fmt.Errorf("%s", stderr)
	}
	r.RepositoryStatus.Cloned = true
	return stdout, nil
}

func CheckoutBranch(r *Repository, branch string) (string, error) {
	cmd := fmt.Sprintf("git checkout -b %s", branch)
	stdout, stderr, err := ShellOut(cmd, r.BaseDir)
	if err != nil {
		log.Printf("Error checking out branch %s: %s", branch, stderr)
		return stderr, err
	}
	r.RepositoryStatus.BranchCreated = true
	return stdout, err
}

type GitHubPrResponse struct {
	CreatedBy []struct {
		Number            int    `json:"number"`
		State             string `json:"state"`
		StatusCheckRollup []any  `json:"statusCheckRollup"`
		URL               string `json:"url"`
	} `json:"createdBy"`
	CurrentBranch struct {
		Number            int    `json:"number"`
		State             string `json:"state"`
		StatusCheckRollup []any  `json:"statusCheckRollup"`
		URL               string `json:"url"`
	} `json:"currentBranch"`
	NeedsReview []any `json:"needsReview"`
}

func NewPullRequestJob(r *Repository, project *Project) *Job {
	prMessage := path.Join(project.ProjectDir, PullRequestFilename)
	cmd := fmt.Sprintf("gh pr create --title '%s' --body-file '%s' --head '%s'", project.PullRequestTitle, prMessage, project.Name)
	job := NewInternalJob(r, cmd)
	job.OnComplete = func(job *Job) {
		job.Repository.RepositoryStatus.PullRequestCreated = true
	}

	return job
}

func NewPushJob(r *Repository, project *Project) *Job {
	cmd := fmt.Sprintf("git push origin '%s'", project.Name)
	return NewInternalJob(r, cmd)
}

func NewPRStatusJob(r *Repository) *Job {
	cmd := "gh pr status --json number,url,state,statusCheckRollup"
	j := NewInternalJob(r, cmd)
	j.OnComplete = func(job *Job) {
		var resp GitHubPrResponse
		err := json.Unmarshal([]byte(strings.Join(job.StdOut, "\n")), &resp)
		if err != nil {
			log.Printf("Error unmarshalling GitHub PR response: %s", err)
			return
		}

		if resp.CurrentBranch.Number == 0 { // No PR for the current branch
			job.Repository.PullRequest = nil
		} else {
			prInfo := PullRequest{
				Number:      resp.CurrentBranch.Number,
				Url:         resp.CurrentBranch.URL,
				Status:      resp.CurrentBranch.State,
				LastChecked: time.Now(),
			}

			job.Repository.PullRequest = &prInfo
			job.Repository.RepositoryStatus.PullRequestCreated = true
		}
	}
	return j
}

func GetRepositoryInfo(repository string) (RepositoryInfo, error) {
	cmd := fmt.Sprintf("gh repo view %s --json name,url,pushedAt", repository)
	var info RepositoryInfo

	stdout, stderr, err := ShellOut(cmd, ".")
	if err != nil {
		return info, fmt.Errorf("%s", stderr)
	}
	err = json.Unmarshal([]byte(stdout), &info)
	if err != nil {
		return info, err
	}
	return info, nil
}
