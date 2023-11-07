package cmd

import (
	"github.com/b-/gomox-uf/cmd/clone"
	"github.com/b-/gomox-uf/cmd/destroy"
	"github.com/b-/gomox-uf/cmd/pveVersion"
	"github.com/b-/gomox-uf/cmd/start"
	"github.com/b-/gomox-uf/cmd/stop"
	"github.com/urfave/cli/v2"
)

func Commands() cli.Commands {
	return cli.Commands{
		start.Command,
		stop.Command,
		pveVersion.Command,
		clone.Command,
		destroy.Command,
	}
}
