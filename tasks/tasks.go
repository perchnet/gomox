package tasks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kr/pretty"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

const (
	NoContent             = "no content" // filter this out from logs
	DefaultTimeout        = 60 * time.Second
	DefaultPollDuration   = time.Millisecond * 500
	DefaultSpinnerCharSet = 9 // Classic Unix quiet |/-\|
	// DefaultSpinnerCharSet = 14 // Docker Compose-style Braille spinner
	spinnerSpeed = 100 * time.Millisecond
)

type spinnerConfig struct {
	enabled       bool
	charSet       int
	randomCharSet bool
	speed         time.Duration
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
	spinnerConfig spinnerConfig
	pollingConfig pollingConfig
}
type pollingConfig struct {
	pollDuration      time.Duration
	timeout           time.Duration
	stopTaskOnTimeout bool
}

type WaitOption func(c *waitConfig)

func WithOutput() WaitOption {
	return func(c *waitConfig) { c.quiet = false }
}

func WithSpinner(opts ...SpinnerOption) WaitOption {
	s := &spinnerConfig{
		charSet: DefaultSpinnerCharSet,
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

// PollingOption is a function used to set items in the pollingConfig struct.
type PollingOption func(c *pollingConfig)

// WithPolling enables polling the task every pollDuration.
// WithPolling also optionally enables a timeout (WithTimeout).
func WithPolling(opts ...PollingOption) WaitOption {
	s := &pollingConfig{
		pollDuration:      DefaultPollDuration,
		timeout:           DefaultTimeout,
		stopTaskOnTimeout: false,
	}
	for _, opt := range opts {
		opt(s)
	}
	return func(c *waitConfig) {
		c.pollingConfig = *s
	}
}

// WithTimeout is a PollingOption that starts a timeout (for `timeout`).
// The timeout will optionally stop the task on timeout (stopTaskOnTimeout)
func WithTimeout(timeout time.Duration, stopTaskOnTimeout bool) PollingOption {
	return func(c *pollingConfig) { c.timeout = timeout; c.stopTaskOnTimeout = stopTaskOnTimeout }
}

// WithPollDuration sets the interval between task polls.
func WithPollDuration(pollDuration time.Duration) PollingOption {
	return func(c *pollingConfig) { c.pollDuration = pollDuration }
}

func Watch(ctx context.Context, start int, t *proxmox.Task) (chan string, error) {
	logrus.Debugf("starting watcher on %s", t.UPID)
	watch := make(chan string)

	log, err := t.Log(ctx, start, 50)
	if err != nil {
		return watch, err
	}

	for i := 0; i < 3; i++ {
		// retry 3 times if the log has no entries
		logrus.Debugf("no logs for %s found, retrying %d of 3 times", t.UPID, i)
		if len(log) > 0 {
			break
		}
		time.Sleep(1 * time.Second)

		log, err = t.Log(ctx, start, 50)
		if err != nil {
			return watch, err
		}
	}

	if len(log) == 0 {
		return watch, fmt.Errorf("no logs available for %s", t.UPID)
	}

	go func() {
		logrus.Debugf("logs found for task %s", t.UPID)
		for _, ln := range log {
			watch <- ln
		}
		logrus.Debugf("watching task %s", t.UPID)
		err := tasktail(ctx, len(log), watch, t)
		if err != nil {
			logrus.Errorf("error watching logs: %s", err)
		}
	}()

	logrus.Debugf("returning watcher for %s", t.UPID)
	return watch, nil
}

func tasktail(ctx context.Context, start int, watch chan string, task *proxmox.Task) error {
	for {
		logrus.Debugf("tailing log for task %s", task.UPID)
		if err := task.Ping(ctx); err != nil {
			return err
		}

		if task.Status != proxmox.TaskRunning {
			logrus.Debugf("task %s is no longer running, closing down watcher", task.UPID)
			close(watch)
			return nil
		}

		logs, err := task.Log(ctx, start, 50)
		if err != nil {
			return err
		}
		for _, ln := range logs {
			watch <- ln
		}
		start = start + len(logs)
		time.Sleep(DefaultPollDuration)
	}
}

func WaitTask(ctx context.Context, task *proxmox.Task, opts ...WaitOption) (err error) {
	taskhead, err := task.Log(ctx, 0, 1)
	if err != nil {
		return err
	}
	taskname := taskhead[0]
	if taskname != NoContent {
		fmt.Println(taskname)
	}
	c := &waitConfig{
		quiet: true, // default to quiet
		spinnerConfig: spinnerConfig{
			enabled: false,
			// charSet: DefaultSpinnerCharSet,
		},
		pollingConfig: pollingConfig{
			pollDuration:      DefaultPollDuration,
			stopTaskOnTimeout: false,
			timeout:           0,
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

	// set up the timeout
	if c.pollingConfig.timeout > 0 { // timeout == 0 means never timeout
		go func() {
			timeout := task.Wait(ctx, c.pollingConfig.pollDuration, c.pollingConfig.timeout)
			proxmox.IsTimeout(timeout)
			if c.pollingConfig.stopTaskOnTimeout {
				err := task.Stop(ctx)
				if err != nil {
					return
				}
			}
		}()
	}

	var (
		msg    string
		newMsg string
	)

	for {
		select {
		/*
			case <-timeoutExpired:
				watch <- taskMsgPrefix(*task, fmt.Sprintln("Timeout expired!"))
				if c.pollingConfig.stopTaskOnTimeout {
					_ = task.Stop(ctx)
				}
				break

		*/

		/*
			case <-ticker.C:
				polledStatus, _ := TaskStatus(ctx, *task)
				if task.IsRunning {
					logrus.Debugln(polledStatus)
				}
				return

		*/

		// receive new task data
		case ln, ok := <-watch:
			if !ok {
				watch = nil
				break
			}
			newMsg = taskMsgPrefix(*task, ln)
			if msg != newMsg {
				msg = newMsg
				// logrus.Infoln(msg)
				s.Suffix = " " + msg
			}
		}
		if watch == nil {
			logrus.Infoln(msg)
			break
		}

	}
	return nil
}

func taskMsgPrefix(task proxmox.Task, msg string) string {
	return fmt.Sprintf("(%s) %s\n", task.Type, msg)
}

// TaskStatus updates the task and returns a message explaining the task's
// status
func TaskStatus(ctx context.Context, task proxmox.Task) (string, error) {
	err := task.Ping(ctx)                               // Update task.
	msg := fmt.Sprintf("%#v\n", pretty.Formatter(task)) // TODO: Improve task status formatting
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
