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
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:        "vmid",
			DefaultText: "",
			FilePath:    "",
			Usage:       "`VMID` to start",
			Required:    true,
			Aliases:     []string{"v"},
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					return fmt.Errorf("VM vmid %d out of range", vmid)
				}
				return (nil)
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

	client := proxmox.NewClient(c.String("pveurl"),
		proxmox.WithCredentials(&credentials),
	)

	cluster, err := client.Cluster(c.Context)
	if err != nil {
		return err
	}

	resources, err := cluster.Resources(c.Context, "vm")
	if err != nil {
		return err
	}

	var vm *proxmox.VirtualMachine
	for _, rs := range resources {
		if rs.VMID == c.Uint64("vmid") {
			node, err := client.Node(c.Context, rs.Node)
			if err != nil {
				return err
			}
			vm, err = node.VirtualMachine(c.Context, int(rs.VMID))
			if err != nil {
				return err
			}
		}
	}

if vm == nil {
		return fmt.Errorf("no vm with id found: %d", c.Uint64("vmid"))
	}

	task, err := vm.Start(context.Background())
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("startvm called! (vm: %#v, task: %#v)", vm, task))
	return nil
}
