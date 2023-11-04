package main

import (
	"io"
	"os"

	"github.com/b-/gomox-uf/cmd"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/urfave/cli/v2"
)
const URL_SUFFIX = "/api2/json"
func main() {
	app := &cli.App{
		Name:     "gomox",
		Usage:    "gomox",
		Commands: cmd.Commands(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "pveuser",
				Aliases: []string{"u"},
				Value: "",
				Usage: "Proxmox VE username",
				EnvVars: []string{"PVE_USER"},
			},
			&cli.StringFlag{
				Name: "pvepassword",
				Aliases: []string{"p"},
				Value: "",
				Usage: "Proxmox VE password",
				EnvVars: []string{"PVE_PASSWORD"},
			},
			&cli.StringFlag{
				Name: "pverealm",
				Aliases: []string{"r"},
				Value: "",
				Usage: "Proxmox VE authentication realm",
				EnvVars: []string{"PVE_REALM"},
			},
			&cli.StringFlag{
				Name: "pveurl",
				Value: "https://127.0.0.1:8006/api2/json",
				Usage: "Proxmox VE API URL",
				EnvVars: []string{"PVE_URL"},
			},
			&cli.StringFlag{
				Name: "scheme",
				Value: "https",
				Usage: "API connection scheme (http or https)",
			},
			&cli.StringFlag{
				Name: "pvehost",
				Aliases: []string{"a"},
				Value: "127.0.0.1",
				Usage: "Proxmox VE hostname/IP address",
				EnvVars: []string{"PVE_HOST"},
			},
			&cli.UintFlag{
				Name: "pveport",
				Value: 8006,
				Usage: "Proxmox VE API port",
				EnvVars: []string{"PVE_PORT"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Turn on verbose debug logging",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Turn on off all logging",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				// treat logrus like fmt.Print
				logrus.SetFormatter(&easy.Formatter{
					LogFormat: "%msg%",
				})
			}
			if ctx.Bool("quiet") {
				logrus.SetOutput(io.Discard)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
