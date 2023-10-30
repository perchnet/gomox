package start

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "start",
	Usage:  "start",
	Action: startVm,
	Flags:  []cli.Flag{
		&cli.Uint64Flag{
			Name:        "vmid",
			DefaultText: "",
			FilePath:    "",
			Usage:       "`VMID` to start",
			Required:    true,
			Aliases:     []string{"c"},
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					return fmt.Errorf("VM vmid %vmid out of range")
				}
				return(nil)
			},
		},
	},
}

func startVm(c *cli.Context) error {
	credentials := proxmox.Credentials{
		Username: c.String("pveuser"),
		Password: c.String("pvepassword"),
		Realm:    c.String("pverealm"),
	}
	vm := proxmox.VirtualMachine{
		VMID: proxmox.StringOrUint64(c.Uint64("vmid")),
	}

	client := proxmox.NewClient(c.String("pveurl"),
		proxmox.WithCredentials(&credentials),
	)

	task, err := vm.Start(context.Background())
	if err != nil {
		panic(err)
	}
	logrus.Info(fmt.Sprintf("startvm called! (vm: %#v, task: %#v)", vm, task))
	return nil
}
