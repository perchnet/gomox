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
	DefaultPollDuration   = time.Second * 5
	DefaultSpinnerCharSet = 9 // Classic Unix quiet |/-\|
	// DefaultSpinnerCharSet = 14 // Docker Compose-style Braille spinner
	spinnerSpeed = 100 * time.Millisecond
)

type spinnerConfig struct {
	enabled bool
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

type waitConfig struct {
	quiet         bool
	timeout       time.Duration
	enablePolling bool // whether to use a ticker to poll the task for completion as well
	pollDuration  time.Duration
	spinnerConfig spinnerConfig
}

type WaitOption func(c *waitConfig)

func WithOutput() WaitOption {
	return func(c *waitConfig) { c.quiet = false }
}

func WithSpinner(opts ...SpinnerOption) WaitOption {
	s := &spinnerConfig{
		charSet: rand.Intn(90), // random spinnerConfig
		speed:   spinnerSpeed,
		enabled: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return func(c *waitConfig) {
		c.quiet = false
		c.spinnerConfig = *s
	}
}

func WithPolling(pollDuration time.Duration, timeout time.Duration) WaitOption {
	return func(c *waitConfig) {
		c.enablePolling = true
		c.pollDuration = pollDuration
		c.timeout = timeout
	}
}

func WaitTask(ctx context.Context, task proxmox.Task, opts ...WaitOption) (err error) {
	c := &waitConfig{
		quiet:         true,  // default to quiet
		timeout:       0,     // No timeout
		enablePolling: false, // No polling
		pollDuration:  DefaultPollDuration,
		spinnerConfig: spinnerConfig{
			enabled: false,
			// charSet: DefaultSpinnerCharSet,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	watch, err := task.Watch(ctx, 0)
	if err != nil {
		return err
	}

	// set up the spinner
	s := spinner.New(spinner.CharSets[c.spinnerConfig.charSet], c.spinnerConfig.speed)
	if c.spinnerConfig.enabled {
		s.Enable()
	} else {
		s.Disable()
	}
	s.Start()
	defer s.Stop()

	// set up the task poller
	if c.enablePolling {
		TaskPoller(ctx, c.pollDuration, c.timeout, task, watch, false)
	}

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
			break
		}
	}
	return nil
}

func TaskPoller(
	ctx context.Context,
	pollDuration time.Duration,
	timeout time.Duration,
	task proxmox.Task,
	watchChannel chan string,
	stopTask bool,
) {
	ticker := time.NewTicker(pollDuration)
	defer ticker.Stop()
	timeoutExpired := make(chan bool)
	if timeout > 0 { // timeout == 0 means never timeout
		go func() {
			time.Sleep(timeout)
			timeoutExpired <- true
			if stopTask {
				err := task.Stop(ctx)
				if err != nil {
					return
				}
			}
		}()
	}
	for {
		select {
		case <-timeoutExpired:
			watchChannel <- fmt.Sprintln("Timed out!")

			close(watchChannel)
			return
		case <-ticker.C:
			polledStatus, _ := TaskStatus(ctx, &task)
			if task.IsRunning {
				logrus.Debugln(polledStatus)
			}
			return
		}
	}
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
