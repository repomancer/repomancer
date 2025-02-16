package internal

import (
	"log"
	"strings"
	"time"
)

// Job represents a shell command run asynchronously in a specific repository's directory
// The command's stdout, stderr and error are captured.
type Job struct {
	Repository *Repository
	Command    string
	Directory  string
	StdOut     string
	StdErr     string
	Error      error
	StartTime  time.Time
	EndTime    time.Time
	Finished   bool
	// This job was created directly by the user (vs commands run internally)
	InternalCommand bool
	// Function run when the job is complete
	OnComplete func(*Job)
}

func NewJob(repository *Repository, command string) *Job {
	return &Job{
		Repository:      repository,
		Command:         command,
		Directory:       repository.BaseDir,
		Finished:        false,
		InternalCommand: false,
	}
}

func NewInternalJob(repository *Repository, command string) *Job {
	job := NewJob(repository, command)
	job.InternalCommand = true
	return job
}

func (j *Job) BuildLogString() string {
	var output []string
	output = append(output, j.Command)
	if strings.TrimSpace(j.StdOut) != "" {
		output = append(output, j.StdOut)
	}
	if strings.TrimSpace(j.StdErr) != "" {
		output = append(output, j.StdErr)
	}
	return strings.Join(output, "\n")
}

func (j *Job) Run() {
	log.Printf("Running command: %s in %s", j.Command, j.Directory)
	j.StartTime = time.Now()
	var err error
	// TODO: Stop storing logs per job, since it's being stored at the repository level???
	j.StdOut, j.StdErr, err = ShellOut(j.Command, j.Directory)
	if err != nil {
		j.Error = err
		// Repository to a function that checks the result of the last job
		j.Repository.LastCommandResult = err
	} else {
		j.Repository.LastCommandResult = nil
	}
	j.Repository.Log(j.BuildLogString())
	j.EndTime = time.Now()
	if j.Error == nil && j.OnComplete != nil {
		j.OnComplete(j)
	}
	j.Finished = true
}

func (j *Job) Duration() string {
	return j.EndTime.Sub(j.StartTime).String()
}
