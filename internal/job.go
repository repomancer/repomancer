package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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
	Output     []byte
	Error      error
	StartTime  time.Time
	EndTime    time.Time
	Finished   bool
	// This job was created directly by the user (vs commands run internally)
	InternalCommand bool
	// Function run when the job is complete
	OnComplete func(*Job)
	LogPath    string
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
	j.StartTime = time.Now()

	var logfile io.Writer
	var err error
	if j.InternalCommand {
		logfile = bytes.NewBuffer([]byte{})
	} else {
		logfile, err = os.OpenFile(j.Repository.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer func(logfile *os.File) {
			err := logfile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(logfile.(*os.File))
		_, err = logfile.Write([]byte(fmt.Sprintf("\n[%s] Running: %s in %s\n",
			time.Now().Format(time.RFC1123),
			j.Command,
			j.Directory)))
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := exec.Command(ShellToUse, "-c", j.Command)
	cmd.Dir = j.Directory
	cmd.Stderr = logfile
	cmd.Stdout = logfile
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	jobErr := cmd.Wait()
	if jobErr != nil {
		j.Error = jobErr
		j.Repository.LastCommandResult = jobErr
	} else {
		j.Repository.LastCommandResult = nil
	}

	j.EndTime = time.Now()
	var finishedMessage string
	if !j.InternalCommand && logfile != nil {
		if jobErr != nil {
			finishedMessage = fmt.Sprintf(
				"[%s] Finished: %s in %s\n",
				time.Now().Format(time.RFC1123),
				jobErr,
				j.EndTime.Sub(j.StartTime))
		} else {
			finishedMessage = fmt.Sprintf(
				"[%s] Finished successfully in %s\n",
				time.Now().Format(time.RFC1123),
				j.EndTime.Sub(j.StartTime))
		}

		_, err = logfile.Write([]byte(finishedMessage))
		if err != nil {
			log.Fatal(err)
		}
	}

	if j.InternalCommand {
		r := bufio.NewReader(logfile.(*bytes.Buffer))
		j.Output, err = io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
	}

	if j.Error == nil && j.OnComplete != nil {
		j.OnComplete(j)
	}
	j.Finished = true
}

func (j *Job) Duration() string {
	return j.EndTime.Sub(j.StartTime).String()
}
