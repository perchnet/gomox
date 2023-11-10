package destroy

import (
	"context"
	"fmt"

	"github.com/b-/gomox/util"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "destroy",
	Usage:  "Delete a virtual machine",
	Action: destroyVmCmd,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:     "vmid",
			Usage:    "`VMID` to delete",
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
			Name:  "force",
			Usage: "If the VM is not stopped, stop before attempting removal.",
		},
		&cli.BoolFlag{
			Name:  "idempotent",
			Usage: "Don't return error if VM is already in requested state",
			Value: false,
		},
	},
}

func destroyVmCmd(c *cli.Context) error {
	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	vmid := c.Uint64("vmid")

	vm, err := util.GetVirtualMachineByVMID(vmid, client, c.Context)
	if err != nil {
		// if we receive an error
		msg := fmt.Sprintf(
			"Could not destroy VM %d.\n"+
				"%#v", vmid, err,
		)
		switch c.Bool("idempotent") {
		case true:
			logrus.Warn(msg)
			return nil
		case false:
			logrus.Panic(msg)
			// don't need to return because the panic will return for us
		}
	}
	if vm.IsStopped() {
		task, err := util.DestroyVm(vm, context.Background())
		if err != nil {
			return err
		}
		logrus.Infof(
			"Deletion requested!\n"+
				"%#v", task,
		)
		if c.Bool("wait") {

		}
		return nil
	} else {
		if c.Bool("force") {
			logrus.Warnf(
				"VM %d is currently %s!\n"+
					"Requesting stop.", vmid, vm.Status,
			)
			task, err := util.RequestState(
				util.StateRequestParams{RequestedState: util.StoppedState, Vm: vm},
				context.Background(),
			)
			if err != nil {
				return err
			} else {
				logrus.Infof(
					"Stop requested!\n"+
						"%#v", task,
				)
			}
		} else {
			err = fmt.Errorf(
				"VM %d is currently %s!\n"+
					"Stop it first, or use `--force`.", vmid, vm.Status,
			)
		}
	}
	return err
}
