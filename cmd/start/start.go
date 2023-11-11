package start

import (
	"fmt"

	"github.com/b-/gomox/util"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "start",
	Usage:  "start a virtual machine",
	Action: startVm,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:     "vmid",
			Usage:    "`VMID` to start",
			Required: true,
			Aliases:  []string{"v"},
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					return fmt.Errorf("VM vmid %d out of range", vmid)
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:  "idempotent",
			Usage: "Don't return error if VM is already in requested state",
			Value: false,
		},
	},
}

// Starts a Proxmox VM as specified by the `vmid` arg
func startVm(c *cli.Context) error {
	requestedState := util.RunningState

	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)

	vmid := c.Uint64("vmid")

	vm, err := util.GetVirtualMachineByVMID(c.Context, vmid, client)
	if err != nil {
		return err
	}

	if vm.IsRunning() {
		msg := fmt.Sprintf("VM %d already in requested state (%s)", vm.VMID, vm.Status)
		switch c.Bool("idempotent") {
		case true:
			logrus.Warn(msg)
			return nil
		case false:
			return fmt.Errorf(msg)
		}
	}
	task, err := util.RequestState(
		c.Context,
		util.StateRequestParams{RequestedState: requestedState, Vm: vm},
	)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("state requested! %#v", task))
	return nil
}
