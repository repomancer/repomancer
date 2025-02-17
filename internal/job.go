package internal

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Job represents a shell command run asynchronously in a specific repository's directory
// The command's stdout, stderr and error are captured.
// It is expected (but not specifically enforced) that only a single Job will be executing in
// a repository at one time, and Jobs will always be executed in the order that they were
// enqueued.
type Job struct {
	Repository *Repository
	Command    string
	Directory  string
	StdOut     []string
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

func (j *Job) Run() {
	log.Printf("Running command: %s in %s", j.Command, j.Directory)
	j.StartTime = time.Now()
	j.Repository.Log(fmt.Sprintf("Command: %s in %s", j.Command, j.Directory))

	cmd := exec.Command(ShellToUse, "-c", j.Command)
	cmd.Dir = j.Directory
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stderr = cmd.Stdout
	// TODO: bufio.NewScanner has a maximum single line length of 4kb and a line longer than
	// that will cause the program to exit with "bufio.Scanner: token too long".
	// This seems acceptable for now, but could be better. This isn't really intended to dump
	// huge amounts of logs in to memory.
	// It's not yet clear what the best behavior would be. Just truncate that line and keep
	// going? Alternately, write all the output to a log file per repository (bonus: durable
	// logs between program runs. Disadvantage: more garbage spread over the filesystem, though
	// that seems relatively minor compared to all the repository cloning. If that was done,
	// replace the Log viewer with just something that would tail the log file. Seems like it
	// would still be nice to have some "intelligence" to it, like highlighting commands.
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		j.Repository.Log(scanner.Text())
		j.StdOut = append(j.StdOut, scanner.Text())
	}
	if scanner.Err() != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		log.Fatal(scanner.Err())
	}
	err = cmd.Wait()
	if err != nil {
		j.Error = err
		j.Repository.LastCommandResult = err
	} else {
		j.Repository.LastCommandResult = nil
	}

	j.EndTime = time.Now()
	if j.Error == nil && j.OnComplete != nil {
		j.OnComplete(j)
	}
	j.Finished = true
}

func (j *Job) Duration() string {
	return j.EndTime.Sub(j.StartTime).String()
}
