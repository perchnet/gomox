package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

const (
	DefaultPollInterval         = time.Duration(time.Second * 5)
	classicSpinnerCharSet       = 9
	dockerComposeSpinnerCharSet = 14
	japaneseSpinnerCharSet      = 25
	oOoSpinnerCharSet           = 28
	spinnerSpeed                = 100 * time.Millisecond
)

// var spinnerType = spinner.CharSets[classicSpinnerCharSet]
var spinnerType = spinner.CharSets[oOoSpinnerCharSet]

// TailTaskStatus waits for a task to complete, displaying a spinner with status messages as it progresses.
func TailTaskStatus(task *proxmox.Task, interval time.Duration, ctx context.Context) error {
	// every interval seconds
	s := spinner.New(spinnerType, spinnerSpeed) // Build our new spinner
	s.Start()                                   // Start the spinner
	taskStatus, err := TaskStatus(*task, ctx)   // init vars, get initial task status
	if err != nil {
		return err
	}
	logrus.Info(taskStatus, "\n")
	s.Suffix = taskStatus // update spinner text
	for {                 // loop
		lastTaskStatus := taskStatus
		taskStatus, err = TaskStatus(*task, ctx) // get new task status
		if err != nil {
			return err
		}
		if taskStatus != lastTaskStatus { // if taskStatus has changed, then
			logrus.Info(taskStatus) // log new taskStatus
			s.Suffix = taskStatus   // update spinner text
		}
		if task.IsRunning {
			time.Sleep(interval)
		} else { // task is not running
			s.Stop()
			break // escape the loop
		}

	}
	return fmt.Errorf("We shouldn't end up here...\n")
}

// TaskStatus updates the task and returns a message explaining the task's status
func TaskStatus(task proxmox.Task, ctx context.Context) (string, error) {
	err := task.Ping(ctx)           // Update task.
	msg := fmt.Sprintf("%#v", task) // TODO: Improve task status formatting
	if task.IsFailed {
		err = fmt.Errorf("the task has failed")
	}
	return msg, err
}

// QuietWaitTask silently waits for a task to complete, without spawning a spinner.
/* Usage example:
err := tasks.QuietWaitTask(
	task,
	tasks.DefaultPollInterval,
	c.Context,
)
if err != nil {
	return err
*/
func QuietWaitTask(task *proxmox.Task, interval time.Duration, ctx context.Context) error {
	taskStatus, err := TaskStatus(*task, ctx) // init vars, get initial task status
	if err != nil {
		return err
	}
	logrus.Tracef("%#v\n", taskStatus)
	for { // loop
		taskStatus, err = TaskStatus(*task, ctx) // get new task status
		if err != nil {
			return err
		}
		if task.IsRunning {
			time.Sleep(interval)
		} else { // task is not running
			break // escape the loop
		}

	}
	return fmt.Errorf("We shouldn't end up here...\n")
}
