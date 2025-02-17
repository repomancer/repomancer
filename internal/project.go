package internal

import (
	"encoding/json"
	"fmt"
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
)

func (p *Project) Select(selectRange SelectRange) {
	cnt := 0
	if selectRange == All {
		for i := 0; i < p.RepositoryCount(); i++ {
			p.GetRepository(i).Selected = true
			cnt++
		}
	} else if selectRange == None {
		for i := 0; i < p.RepositoryCount(); i++ {
			p.GetRepository(i).Selected = false
		}
	} else if selectRange == Errors {
		for i := 0; i < p.RepositoryCount(); i++ {
			if p.GetRepository(i).LastCommandResult != nil {
				p.GetRepository(i).Selected = true
				cnt++
			} else {
				p.GetRepository(i).Selected = false
			}
		}
	} else if selectRange == TenMore {
		added := 0
		for i := 0; i < p.RepositoryCount(); i++ {
			if !p.GetRepository(i).Selected {
				if added < 10 {
					p.GetRepository(i).Selected = true
					added++
					cnt++
				}
			}
		}
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

type Project struct {
	mu           sync.Mutex
	Name         string
	Description  string
	Repositories []*Repository
	ProjectDir   string `json:"-"` // Calculated on load, not saved with configuration
}

func (p *Project) AddRepository(host, org, name string) error {
	r := &Repository{
		Host:         host,
		Organization: org,
		Name:         name,
		BaseDir:      path.Join(p.ProjectDir, host, org, name),
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
	p.mu.Lock()
	p.Repositories = append(p.Repositories, r)
	p.mu.Unlock()
	_, err = Clone(r)
	if err != nil {
		return err
	}
	_, err = CheckoutBranch(r, p.Name)

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
		cnt += p.GetRepository(i).JobCount()
	}
	return cnt
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
	}

	return &payload, nil
}

func OpenProject(projectPath string) (*Project, error) {
	project, err := ReadProjectConfig(projectPath)
	if err != nil {
		return &Project{}, err
	} else {
		return project, nil
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
				Name:        name,
				Description: description,
				ProjectDir:  projectPath,
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
			return &project, nil
		} else {
			return &Project{}, err
		}
	}

	return &Project{}, err
}
