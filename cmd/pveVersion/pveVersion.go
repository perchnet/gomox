package pveVersion

import (
	"github.com/b-/gomox/util"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "pveVersion",
	Usage:  "pveVersion",
	Action: pveVersion,
	Flags:  []cli.Flag{},
}

func pveVersion(c *cli.Context) error {
	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)

	version, err := client.Version(c.Context)
	if err != nil {
		return err
	}

	logrus.Info(version.Release)
	return nil
}
