package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"strings"
	"testing"
)

func NewTestRepository() *Repository {
	r := Repository{}
	return &r
}

func TestJobCommandFailed(t *testing.T) {
	r := NewTestRepository()
	j := NewInternalJob(r, "ThisDoesNotExist")
	j.Run()
	assert.Equal(t, j.Finished, true)
	assert.NotNil(t, j.Error)
	assert.True(t, strings.HasSuffix(j.Duration(), "ms"))
}

func TestNewInternalJob(t *testing.T) {
	r := NewTestRepository()
	j := NewInternalJob(r, "ls")
	assert.Equal(t, j.Finished, false)
	assert.Equal(t, j.InternalCommand, true)
}

func TestRunningJob(t *testing.T) {
	r := NewTestRepository()
	tmpDir := t.TempDir()
	r.BaseDir = tmpDir
	r.LogFile = path.Join(tmpDir, "log.txt")
	j := NewJob(r, "pwd")
	called := false
	j.OnComplete = func(j *Job) {
		called = true
	}

	j.Run()
	assert.Equal(t, j.Finished, true)
	assert.Equal(t, j.InternalCommand, false)

	b, err := os.ReadFile(r.LogFile)
	assert.NoError(t, err)
	s := string(b)
	assert.Contains(t, s, "Running: pwd in")
	assert.Contains(t, s, "Finished successfully in")

	assert.True(t, called)
}
