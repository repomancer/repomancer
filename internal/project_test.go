package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const ConfigFile = `
{
  "Name": "Test",
  "PullRequestDescription": "Testing",
  "Repositories": [
    {"Name": "github.com/jashort/foo"},
    {"Name": "github.com/jashort/bar"}
  ],
  "ProjectDir": "/var/folders/wm/h0mhbkg91lj2sfcr_g9z5p3h0000gn/T/TestCreateOrOpenProject1573936796/001"
}
`

func TestCreateProject(t *testing.T) {
	tmpDir := t.TempDir()

	project, _ := CreateProject("Test", "Testing", tmpDir)

	if project.ProjectDir != tmpDir {
		t.Errorf("project dir does not match")
	}
}

func TestCreateProject_Fail(t *testing.T) {
	_, err := CreateProject("Test", "Testing", "/cannotwrite")
	if err == nil {
		t.Errorf("error should not be nil")
	}
}

func TestOpenExistingProject(t *testing.T) {
	tmpDir := t.TempDir()
	f, _ := os.Create(filepath.Join(tmpDir, "config.json"))
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	_, _ = f.Write([]byte(ConfigFile))

	project, err := OpenProject(tmpDir)

	if err != nil {
		t.Errorf("error should be nil, was %v", err)
	} else {
		if project.ProjectDir != tmpDir {
			t.Errorf("project dir does not match")
		}
		if project.Name != "Test" {
			t.Errorf("project name does not match")
		}
		if project.PullRequestDescription != "Testing" {
			t.Errorf("project description does not match")
		}
	}
}

func TestProjectRepositories(t *testing.T) {
	tmpDir := t.TempDir()
	project, _ := CreateProject("Test", "Testing", tmpDir)
	for i := range 3 {
		project.Repositories = append(project.Repositories, &Repository{
			Name: fmt.Sprintf("github.com/jashort/test%d", i),
		})
	}

	assert.Equal(t, 3, project.RepositoryCount())
	// If no repositories are selected, all repositories are selected
	assert.Equal(t, 3, len(project.SelectedRepositories()))
	project.Select(All)
	assert.Equal(t, 3, len(project.SelectedRepositories()))
	// Fake an open pull request and make sure it is the only one selected
	project.Repositories[0].PullRequest = &PullRequest{
		Number:      1,
		LastChecked: time.Time{},
	}
	project.Select(ReposWithPullRequest)
	assert.Equal(t, 1, len(project.SelectedRepositories()))
	// Delete the selected repository
	project.DeleteSelectedRepositories()
	assert.Equal(t, 2, len(project.SelectedRepositories()))
	assert.Equal(t, 2, project.RepositoryCount())
}
