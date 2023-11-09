package cmd

import (
	"github.com/b-/gomox-uf/cmd/cloneCmd"
	"github.com/b-/gomox-uf/cmd/destroyCmd"
	"github.com/b-/gomox-uf/cmd/pveVersionCmd"
	"github.com/b-/gomox-uf/cmd/startCmd"
	"github.com/b-/gomox-uf/cmd/stopCmd"
	"github.com/b-/gomox-uf/cmd/taskStatusCmd"
	"github.com/urfave/cli/v2"
)

func Commands() cli.Commands {
	return cli.Commands{
		startCmd.Command,
		stopCmd.Command,
		pveVersionCmd.Command,
		cloneCmd.Command,
		destroyCmd.Command,
		taskStatusCmd.Command,
	}
}
