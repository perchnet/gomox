package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

const (
	DefaultPollInterval = time.Second * 5
	// SpinnerCharSet       = 9 // Classic Unix quiet |/-\|
	// SpinnerCharSet = 14 // Docker Compose-style Braille quiet
	spinnerCharSet = 25 // Japanese quiet
	// SpinnerCharSet = 28 // oOo
	spinnerSpeed = 100 * time.Millisecond
)

type spinnerConfig struct {
	charSet int
	speed   time.Duration
}

type SpinnerOption func(c *spinnerConfig)

func WithSpinnerCharSet(charset int) SpinnerOption {
	return func(c *spinnerConfig) { c.charSet = charset }
}
func WithSpinnerSpeed(speed time.Duration) SpinnerOption {
	return func(c *spinnerConfig) { c.speed = speed }
}

// NewSpinner returns a *Spinner with the specified options
func NewSpinner(opts ...SpinnerOption) *spinner.Spinner {
	c := &spinnerConfig{
		charSet: rand.Intn(90), // random spinner
		speed:   100 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(c)
	}
	return spinner.New(spinner.CharSets[c.charSet], c.speed)
}

type waitConfig struct {
	quiet   bool
	spinner spinnerConfig
}

type WaitOption func(c *waitConfig)

func WithOutput(quiet bool, spinnerOpt ...SpinnerOption) WaitOption {
	return func(c *waitConfig) { c.quiet = false }
}

func WithSpinner(opts ...SpinnerOption) WaitOption {
	s := &spinnerConfig{
		charSet: rand.Intn(90), // random spinner
		speed:   100 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(s)
	}
	return func(c *waitConfig) {
		c.quiet = false
		c.spinner = *s
	}
}

func WaitTask(ctx context.Context, task proxmox.Task, opts ...WaitOption) (err error) {
	c := &waitConfig{
		quiet: true, // default to quiet
	}
	for _, opt := range opts {
		opt(c)
	}
	watch, err := task.Watch(ctx, 0)
	if err != nil {
		return err
	}
	s := spinner.New(spinner.CharSets[c.spinner.charSet], c.spinner.speed)
	s.Start()
	var newMsg string
	for {
		select {
		case ln, ok := <-watch:
			if !ok {
				watch = nil
				return err
			}
			newMsg = taskMsgPrefix(task, ln)
			if s.Suffix != " "+newMsg {
				switch c.quiet {
				case true:
					logrus.Trace(s.Suffix + "\n")
				case false:
					logrus.Info(s.Suffix + "\n")
				}
			}
			s.Suffix = " " + // looks better with a space…
				newMsg
		}
		if watch == nil || !(task.IsRunning) {
			s.Stop()
			break
		}
	}
	return nil
}

func taskMsgPrefix(task proxmox.Task, msg string) string {
	return fmt.Sprintf("(%s) %s", task.Type, msg)
}

// TaskStatus updates the task and returns a message explaining the task's
// status
func TaskStatus(ctx context.Context, task *proxmox.Task) (string, error) {
	err := task.Ping(ctx)            // Update task.
	msg := fmt.Sprintf("%#v", *task) // TODO: Improve task status formatting
	if task.IsFailed {
		err = fmt.Errorf("the task has failed")
	}
	return msg, err
}

func GetWaitCmd(task proxmox.Task) string {
	return fmt.Sprintf(
		`
To watch the running operation, run:
    %s -w taskstatus "%s"
`, os.Args[0], task.UPID,
	)
}
