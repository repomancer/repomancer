package internal

import (
	"github.com/stretchr/testify/assert"
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
