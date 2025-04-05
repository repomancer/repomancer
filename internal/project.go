package internal

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type SelectRange int

const (
	All SelectRange = iota
	None
	Errors
	TenMore
	ReposWithPullRequest
	ReposWithoutPullRequest
)

const RepositoryQueueSize = 512 // Maximum number of queued repositories to work on TODO: real queue instead of channels

type Project struct {
	mu                     sync.Mutex
	Name                   string
	PullRequestTitle       string
	PullRequestDescription string
	Repositories           []*Repository
	ProjectDir             string           `json:"-"` // Calculated on load, not saved with configuration
	WorkerChannel          chan *Repository `json:"-"`
}

func (p *Project) Select(selectRange SelectRange) {
	switch selectRange {
	case All:
		for i := 0; i < p.RepositoryCount(); i++ {
			p.GetRepository(i).Selected = true
		}
	case None:
		for i := 0; i < p.RepositoryCount(); i++ {
			p.GetRepository(i).Selected = false
		}
	case Errors:
		for i := 0; i < p.RepositoryCount(); i++ {
			if p.GetRepository(i).LastCommandResult != nil {
				p.GetRepository(i).Selected = true
			} else {
				p.GetRepository(i).Selected = false
			}
		}
	case TenMore:
		added := 0
		for i := 0; i < p.RepositoryCount(); i++ {
			if !p.GetRepository(i).Selected {
				if added < 10 {
					p.GetRepository(i).Selected = true
					added++

				}
			}
		}
	case ReposWithPullRequest:
		for i := 0; i < p.RepositoryCount(); i++ {
			if p.GetRepository(i).PullRequest != nil {
				p.GetRepository(i).Selected = true
			} else {
				p.GetRepository(i).Selected = false
			}
		}

	case ReposWithoutPullRequest:
		for i := 0; i < p.RepositoryCount(); i++ {
			if p.GetRepository(i).PullRequest == nil {
				p.GetRepository(i).Selected = true
			} else {
				p.GetRepository(i).Selected = false
			}
		}

	default:
		panic("unhandled default case")
	}
}

// SelectedRepositories will return only selected repositories if anything is selected, or all repositories
// if nothing is selected
func (p *Project) SelectedRepositories() []*Repository {
	p.mu.Lock()
	defer p.mu.Unlock()
	var repos []*Repository

	anySelected := false
	for _, repo := range p.Repositories {
		if repo.Selected {
			anySelected = true
			break
		}
	}

	if anySelected {
		for i := 0; i < len(p.Repositories); i++ {
			if p.Repositories[i].Selected {
				repos = append(repos, p.Repositories[i])
			}
		}
	} else {
		repos = append(repos, p.Repositories...)
	}
	return repos
}

func (p *Project) AddJobToRepositories(cmd string) {
	selected := p.SelectedRepositories()
	for i := 0; i < len(selected); i++ {
		j := NewJob(selected[i], cmd)
		selected[i].AddJob(j)
	}
}

func (p *Project) AddInternalJobToRepositories(cmd string, onComplete func(job *Job)) {
	selected := p.SelectedRepositories()
	for i := 0; i < len(selected); i++ {
		j := NewInternalJob(selected[i], cmd)
		j.OnComplete = onComplete
		selected[i].AddJob(j)
	}
}

func (p *Project) AddRepository(host, org, name string) error {
	log.Printf("Adding repository %s/%s/%s", host, org, name)
	r := &Repository{
		Host:         host,
		Organization: org,
		Name:         name,
		BaseDir:      path.Join(p.ProjectDir, host, org, name),
		LogFile:      path.Join(p.ProjectDir, fmt.Sprintf("%s_%s_%s.log", host, org, name)),
		Jobs:         NewJobQueue(),
		RepositoryStatus: RepositoryStatus{
			Cloned:             false,
			BranchCreated:      false,
			PullRequestCreated: false,
			PullRequestClosed:  false,
		},
	}

	if unix.Access(r.BaseDir, unix.W_OK) == nil {
		return fmt.Errorf("%s already exists", r.BaseDir)
	}

	err := os.MkdirAll(r.BaseDir, 0755)
	if err != nil {
		return err
	}
	_, err = Clone(r)
	if err != nil {
		// If there was an error cloning, remove the directory that was created
		_ = os.RemoveAll(r.BaseDir)
		return err
	}
	_, err = CheckoutBranch(r, p.Name)

	p.mu.Lock()
	p.Repositories = append(p.Repositories, r)
	p.mu.Unlock()
	log.Printf("Cloned repository %s to %s", r.Name, r.BaseDir)

	err = p.SaveProject()
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) SaveProject() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	data, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(p.ProjectDir, "config.json"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) AddRepositoryFromUrl(url string) error {
	normalized, err := NormalizeGitUrl(url)
	if err != nil {
		return err
	}
	s := strings.Split(normalized, "/")
	if len(s) != 3 {
		return fmt.Errorf("invalid repository url: %s (expecting host/org/name", url)
	}

	return p.AddRepository(s[0], s[1], s[2])
}

func (p *Project) GetRepository(i int) *Repository {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Repositories[i]
}

func (p *Project) RepositoryCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.Repositories)
}

func (p *Project) SelectedRepositoryCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	cnt := 0
	for i := 0; i < len(p.Repositories); i++ {
		if p.Repositories[i].Selected {
			cnt++
		}
	}
	return cnt
}

func (p *Project) TotalJobCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	cnt := 0
	for i := 0; i < len(p.Repositories); i++ {
		cnt += p.GetRepository(i).Jobs.Len()
	}
	return cnt
}

// DeleteSelectedRepositories removes selected repositories from the project
// and deletes the files associated with them. It does not remove anything from
// the remote repository that has already been pushed.
// If no repositories are selected, all repositories are deleted. This should only
// after giving the user an "are you sure" prompt.
func (p *Project) DeleteSelectedRepositories() {
	selectedCount := p.SelectedRepositoryCount()
	p.mu.Lock()
	defer p.mu.Unlock()

	var toDelete []*Repository
	var toKeep []*Repository

	if selectedCount == 0 {
		toDelete = append(toDelete, p.Repositories...)
	} else {
		for i := 0; i < len(p.Repositories); i++ {
			if p.Repositories[i].Selected {
				toDelete = append(toDelete, p.Repositories[i])
			} else {
				toKeep = append(toKeep, p.Repositories[i])
			}
		}
	}

	p.Repositories = toKeep

	go func() {
		for _, repo := range toDelete {
			err := os.RemoveAll(repo.BaseDir)
			if err != nil {
				log.Printf("Failed to remove repository %s: %s", repo.Name, err)
			}
		}
	}()
	log.Printf("Deleting selected repositories")
}

func (p *Project) DeleteSelectedLogs() {
	toDelete := p.SelectedRepositories()
	p.mu.Lock()
	defer p.mu.Unlock()

	go func() {
		for _, repo := range toDelete {
			err := os.Remove(repo.LogFile)
			if err != nil {
				log.Printf("Failed to remove logfile for %s: %s", repo.Name, err)
			}
		}
	}()
	log.Printf("Deleting selected repository logs")
}

func ReadProjectConfig(projectPath string) (*Project, error) {
	content, err := os.ReadFile(filepath.Join(projectPath, "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
		}
		return &Project{}, err
	}

	var payload Project
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return &Project{}, err
	}
	payload.ProjectDir = projectPath

	for i := 0; i < len(payload.Repositories); i++ {
		repo := payload.Repositories[i]
		repo.BaseDir = path.Join(projectPath, repo.Host, repo.Organization, repo.Name)
		repo.LogFile = path.Join(projectPath, fmt.Sprintf("%s_%s_%s.log", repo.Host, repo.Organization, repo.Name))
		repo.Jobs = NewJobQueue()
	}

	return &payload, nil
}

func OpenProject(projectPath string) (*Project, error) {
	project, err := ReadProjectConfig(projectPath)
	if err != nil {
		return &Project{}, err
	} else {
		project.WorkerChannel = make(chan *Repository, RepositoryQueueSize)
		for w := 1; w <= workerCount(); w++ {
			go worker(w, project.WorkerChannel)
		}
		return project, nil
	}
}

func worker(id int, repositories <-chan *Repository) {
	for repo := range repositories {
		log.Printf("Worker %d started repository %s", id, repo.Name)
		repo.RunJobs()
		log.Printf("Worker %d finished repository %s", id, repo.Name)
	}
}

// CreateProject creates the directory and file if it doesn't exist.
func CreateProject(name, description, projectPath string) (*Project, error) {
	fileInfo, err := os.Stat(projectPath)
	if err == nil && !fileInfo.IsDir() {
		return &Project{}, fmt.Errorf("%s exists but is not a directory", projectPath)
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(projectPath, os.ModePerm)
		if err != nil {
			return &Project{}, err
		}
	}
	_, err = ReadProjectConfig(projectPath)
	if err == nil {
		return &Project{}, fmt.Errorf("project already exists in %s", projectPath)
	} else {
		if os.IsNotExist(err) {
			project := Project{
				Name:                   name,
				PullRequestDescription: description,
				ProjectDir:             projectPath,
				WorkerChannel:          make(chan *Repository, RepositoryQueueSize),
			}

			f, _ := os.Create(filepath.Join(projectPath, "config.json"))
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(f)

			asJson, _ := json.MarshalIndent(&project, "", "  ")
			_, err := f.Write(asJson)
			if err != nil {
				log.Fatal(err)
			}

			for w := 1; w <= workerCount(); w++ {
				go worker(w, project.WorkerChannel)
			}
			return &project, nil
		} else {
			return &Project{}, err
		}
	}
}

// workerCount returns the number of workers to use. This is configurable, but defaults
// to 1 if the Fyne framework is not running.
func workerCount() int {
	currentApp := fyne.CurrentApp()
	if currentApp != nil {
		return currentApp.Preferences().IntWithFallback("workerCount", 5)
	} else {
		return 1
	}
}
