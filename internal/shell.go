package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// RunCommand is for running a shell command that should NOT be run per repository. For that,
// create a Job and use Repository.AddJob()
func RunCommand(dir string, timeoutSeconds int, command string, args ...string) (stdout string, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Create the command with the context
	cmd := exec.CommandContext(ctx, command, args...)

	// Set the working directory if specified
	if dir != "" {
		cmd.Dir = dir
	}

	// Create buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	err = cmd.Run()

	// Check if the context deadline was exceeded
	if errors.Is(context.DeadlineExceeded, ctx.Err()) {
		return "", "", fmt.Errorf("command timed out after %d seconds", timeoutSeconds)
	}

	// Return the captured output and any error
	return stdoutBuf.String(), stderrBuf.String(), err
}
