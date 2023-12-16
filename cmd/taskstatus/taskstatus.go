package taskstatus

import (
	"fmt"

	"github.com/perchnet/gomox/tasks"
	"github.com/perchnet/gomox/util"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const UsageText = "gomox [-w] taskstatus <UPID>"

var Command = &cli.Command{
	Name:      "taskstatus",
	Usage:     "Get the status of a given task, by UPID",
	UsageText: UsageText,
	Action:    taskStatusCmd,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "timeout",
			Category: "wait",
			Usage:    "Wait up to `TIMEOUT` seconds for task completion",
			// Value:    30,
			Aliases: []string{"s"},
		},
		&cli.IntFlag{
			Name:     "interval",
			Category: "wait",
			Usage:    "Update every `INTERVAL` seconds.",
			Value:    1,
			Aliases:  []string{"i"},
		},
	},
}

func taskStatusCmd(c *cli.Context) error {
	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	if len(c.Args().Slice()) == 0 {
		return fmt.Errorf("Usage: " + UsageText)
	}
	upid := c.Args().First()
	tailMode := c.Bool("wait")
	task := proxmox.NewTask(proxmox.UPID(upid), &client)

	taskStatus, err := tasks.TaskStatus(c.Context, *task)
	if err != nil {
		return err
	}
	logrus.Info(taskStatus)
	if task.IsRunning && tailMode {
		err = WaitForCliTask(c, task)
		if err != nil {
			return err
		}
	}

	return nil
}

// WaitForCliTask waits for `task` to complete
func WaitForCliTask(c *cli.Context, task *proxmox.Task) error {
	var err error
	if c.Bool("quiet") {
		err = tasks.WaitTask(
			c.Context,
			task,
		)
		if err != nil {
			return err
		}
	} else {
		err = tasks.WaitTask(
			c.Context,
			task,
			tasks.WithOutput(),
			tasks.WithPolling(),
			tasks.WithSpinner(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
