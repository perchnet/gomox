package cmd

import (
	"github.com/perchnet/gomox/cmd/clone"
	"github.com/perchnet/gomox/cmd/config"
	"github.com/perchnet/gomox/cmd/destroy"
	"github.com/perchnet/gomox/cmd/list"
	"github.com/perchnet/gomox/cmd/pveVersion"
	"github.com/perchnet/gomox/cmd/set"
	"github.com/perchnet/gomox/cmd/start"
	"github.com/perchnet/gomox/cmd/stop"
	"github.com/perchnet/gomox/cmd/taskstatus"
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
		set.Command,
	}
}
