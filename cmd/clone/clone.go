package clone

import (
	"context"
	"fmt"

	"github.com/b-/gomox-uf/internal"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "clone",
	Usage:  "Clone a virtual machine",
	Action: cloneVm,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:     "vmid",
			Usage:    "`VMID` to clone from",
			Required: true,
			Aliases:  []string{"v"},
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					return fmt.Errorf("vmid %d out of range", vmid)
				}
				return nil
			},
		},
		&cli.Uint64Flag{
			Name:        "newid",
			Usage:       "`VMID` for the clone",
			Required:    false,
			Aliases:     []string{"n"},
			DefaultText: "next available",
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					if vmid != 0 {
						return fmt.Errorf("vmid %d out of range", vmid)
					}
				}
				return nil
			},
		},
		&cli.Uint64Flag{
			Name:        "bwlimit",
			Usage:       "Override I/O bandwidth limit (in KiB/s).",
			DefaultText: "unlimited",
		},
		&cli.StringFlag{
			Name:  "description",
			Usage: "Description for the new vm.",
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "Target format for file storage. Only valid for full clone. Can be raw, qcow, or vmdk.",
		},
		&cli.BoolFlag{
			Name:  "full",
			Usage: "Create a full copy of all disks. This is always done when you clone a normal VM. For VM templates, we try to create a linked clone by default.",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "Set a name for the new VM.",
		},
		&cli.StringFlag{
			Name:  "pool",
			Usage: "Add the new VM to the specified pool.",
		},
		&cli.StringFlag{
			Name:  "snapname",
			Usage: "The name of the snapshot.",
		},
		&cli.StringFlag{
			Name:  "storage",
			Usage: "Target storage for full clone.",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: "Target node. Only allowed if the original VM is on shared storage.",
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite the target VMID if it already exists. (Note: only relevant when manually specifying VMID.)",
		},
	},
}

func bool2uint(b bool) uint {
	if b {
		return 1
	}
	return 0
}

func bool2uint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// Clones a Proxmox VM as specified by the `from` arg and `params` struct
func cloneVm(c *cli.Context) error {
	cloneOptions := proxmox.VirtualMachineCloneOptions{
		NewID:    int(c.Uint64("newid")),
		BWLimit:  c.Uint64("bwlimit"),
		Full:     bool2uint8(c.Bool("full")),
		Name:     c.String("name"),
		Pool:     c.String("pool"),
		SnapName: c.String("snapname"),
		Storage:  c.String("storage"),
		Target:   c.String("target"),
	}

	client := internal.InstantiateClient(
		internal.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	vmid := c.Uint64("vmid")
	newid := c.Uint64("newid")
	vm, err := internal.GetVirtualMachineByVMID(vmid, client, c.Context)
	if err != nil {
		return err
	}

	if newid != 0 { // if we're manually assigning the target VMID
		// check if VM already exists with target VMID
		vm, _ := internal.GetVirtualMachineByVMID(newid, client, c.Context)
		if vm != nil {
			logrus.Infof("Virtual machine with target ID %d already exists.\n", newid)
			switch c.Bool("overwrite") {
			case true:
				logrus.Infof("Overwrite requested.\n")
				task, err := internal.DestroyVm(vm, context.Background())
				logrus.Infof("Overwrite requested. Destroying VM %#v.\n%#v", vm, task)
				// need to wait for completion...
				if err != nil {
					logrus.Panicf("Shit!")
				}
			case false:
				logrus.Panicf("Use --overwrite if necessary.\n")
			}

		}
	}

	outvmid, task, err := vm.Clone(context.Background(), &cloneOptions)
	if err != nil {
		logrus.Panic("Oh no! ", err)
	}
	logrus.Infof("clone requested! new id: %d.\n%#v\n", outvmid, task)
	return nil
}
