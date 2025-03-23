package internal

import (
	"bytes"
	"os/exec"
)

const ShellToUse = "bash"

// ShellOut is for running a shell command that should NOT be run per repository. For that,
// create a Job and use Repository.AddJob()
func ShellOut(command string, directory string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Dir = directory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
