package pveVersion

import (
	"context"

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
		Credentials := proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		}
		client := proxmox.NewClient(c.String("pveurl"),
			proxmox.WithCredentials(&Credentials),
		)

		version, err := client.Version(context.Background())
		if err != nil {
			panic(err)
		}

	logrus.Info(version.Release)
	return nil
}
