package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCommand(t *testing.T) {
	gotStdout, _, err := RunCommand("", 5, "ls -l")
	assert.Contains(t, gotStdout, "total")
	assert.NoError(t, err)
}

func TestRunTimeout(t *testing.T) {
	_, _, err := RunCommand("", 1, "sleep 2")
	// If the command times out, stderr could be empty but err will not be nil
	assert.Error(t, err)
}

func TestRunCommandErrors(t *testing.T) {
	_, gotStderr, err := RunCommand("", 1, "cat /doesnotexist")
	// If the command throws an error, stderr will contain the error
	// and err will not be nil
	assert.Contains(t, gotStderr, "cat: /doesnotexist: No such file or directory")
	assert.Error(t, err)
}
