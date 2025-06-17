package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"os/exec"
	"strings"
	"time"
)

// RunCommand is for running a shell command that should NOT be run per repository. For that,
// create a Job and use Repository.AddJob()
func RunCommand(dir string, timeoutSeconds int, command string) (stdout string, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Create the command with the context
	// Run commands inside ZSH with
	args := ShellArgs()
	args = append(args, command)
	cmd := exec.CommandContext(ctx, ShellToUse(), args...)

	// Set the working directory if specified
	if dir != "" {
		cmd.Dir = dir
	}

	// Create buffers for stdin, stdout, stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
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

const DefaultShell = "zsh"
const DefaultShellArgs = "--login -i -c"

func ShellToUse() string {
	if fyne.CurrentApp() == nil {
		return DefaultShell
	}
	return fyne.CurrentApp().Preferences().StringWithFallback("shell", DefaultShell)
}

func ShellArgs() []string {
	var args string
	if fyne.CurrentApp() == nil {
		args = DefaultShellArgs
	} else {
		args = fyne.CurrentApp().Preferences().StringWithFallback("shellArguments", DefaultShellArgs)
	}
	return strings.Split(args, " ")
}
