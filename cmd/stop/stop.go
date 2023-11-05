package stop

import (
	"context"
	"fmt"

	"github.com/b-/gomox-uf/internal/clientinstantiator"
	"github.com/b-/gomox-uf/internal/pveurl"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "stop",
	Usage:  "Stop a virtual machine",
	Action: stopVm,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:        "vmid",
			DefaultText: "",
			FilePath:    "",
			Usage:       "`VMID` to stop",
			Required:    true,
			Aliases:     []string{"v"},
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					return fmt.Errorf("VM vmid %d out of range", vmid)
				}
				return nil
			},
		},
	},
}

func stopVm(c *cli.Context) error {
	client := clientinstantiator.InstantiateClient(
		pveurl.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
	},
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

	if vm.Status == "stopped" {
		msg := fmt.Sprintf("VM %d already in requested state (%s)", vm.VMID, vm.Status)
		switch c.Bool("idempotent") {
		case true:
			logrus.Warn(msg)
			return nil
		case false:
			return fmt.Errorf(msg)
		}
	}

	task, err := vm.Stop(context.Background())
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("stopvm called! (vm: %#v, task: %#v)", vm, task))
	return nil
}
