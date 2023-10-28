package example

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "example",
	Usage:  "",
	Action: stub,
	Flags:  []cli.Flag{},
}

func stub(c *cli.Context) error {
	logrus.Info("example command")
	return nil
}
