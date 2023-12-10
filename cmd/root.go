package cmd

import (
	"github.com/b-/gomox/cmd/clone"
	"github.com/b-/gomox/cmd/config"
	"github.com/b-/gomox/cmd/destroy"
	"github.com/b-/gomox/cmd/list"
	"github.com/b-/gomox/cmd/pveVersion"
	"github.com/b-/gomox/cmd/start"
	"github.com/b-/gomox/cmd/stop"
	"github.com/b-/gomox/cmd/taskstatus"
	"github.com/urfave/cli/v2"
)

func Commands() cli.Commands {
	return cli.Commands{
		start.Command,
		stop.Command,
		pveVersion.Command,
		clone.Command,
		destroy.Command,
		taskstatus.Command,
		list.Command,
		config.Command,
	}
}
