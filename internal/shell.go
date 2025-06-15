package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RunCommand is for running a shell command that should NOT be run per repository. For that,
// create a Job and use Repository.AddJob()
func RunCommand(dir string, timeoutSeconds int, command string, args ...string) (stdout string, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Create the command with the context
	// Run commands inside ZSH with
	var myArgs []string
	myArgs = append(myArgs, "--login")
	cmd := exec.CommandContext(ctx, ShellToUse, myArgs...)

	// Set the working directory if specified
	if dir != "" {
		cmd.Dir = dir
	}

	// Create buffers for stdin, stdout, stderr
	var stdoutBuf, stderrBuf, stdinBuf bytes.Buffer
	stdinBuf.WriteString(command)
	if !strings.HasSuffix(command, "\n") {
		stdinBuf.WriteString("\n")
	}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Stdin = &stdinBuf
	// Run the command
	err = cmd.Run()

	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), err
	}

	// Check if the context deadline was exceeded
	if errors.Is(context.DeadlineExceeded, ctx.Err()) {
		return "", "", fmt.Errorf("command timed out after %d seconds", timeoutSeconds)
	}

	// Return the captured output and any error
	return stdoutBuf.String(), stderrBuf.String(), err
}
