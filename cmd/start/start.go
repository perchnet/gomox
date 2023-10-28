package start

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "start",
	Usage:  "start",
	Action: stub,
	Flags:  []cli.Flag{},
}

func stub(c *cli.Context) error {
	logrus.Info("start!")
	return nil
}
