package pveVersionCmd

import (
	"context"

	"github.com/b-/gomox-uf/internal"
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
	client := internal.InstantiateClient(
		internal.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)

	version, err := client.Version(context.Background())
	if err != nil {
		return err
	}

	logrus.Info(version.Release)
	return nil
}
