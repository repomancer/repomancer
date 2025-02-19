package internal

import (
	"os"
	"path/filepath"
	"testing"
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
