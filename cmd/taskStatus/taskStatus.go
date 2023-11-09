package taskStatus

import (
	"time"

	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const taskMsgFmt = "task: %#v\n"

var Command = &cli.Command{
	Name:   "taskstatus",
	Usage:  "Get the status of a given task, by UPID",
	Action: taskStatusCmd,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "upid",
			Usage:    "Proxmox task ID (`UPID`) to get the status for.",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "timeout",
			Category: "wait",
			Usage:    "Wait up to `TIMEOUT` seconds for task completion",
			Value:    30,
			Aliases:  []string{"s"},
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
	client := InstantiateClient(
		GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	tailMode := c.Bool("wait")
	task := proxmox.NewTask(proxmox.UPID(c.String("upid")), &client)

	taskStatus, err := TaskStatus(*task, c.Context)
	if err != nil {
		return err
	}
	if task.IsRunning {
		if tailMode {
			TailTaskStatus(
				task,
				time.Duration(c.Int("interval"))*time.Second,
				c.Context,
			)
		}
	} else {
		logrus.Info(taskStatus)
	}

	return nil
}
