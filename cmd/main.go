package cmd

import (
	"github.com/b-/gomox-uf/cmd/pveVersion"
	"github.com/b-/gomox-uf/cmd/start"
	"github.com/urfave/cli/v2"
)

func Commands() cli.Commands {
	return cli.Commands{
		start.Command,
		pveVersion.Command,
	}
}
