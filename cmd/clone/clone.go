package clone

import (
	"context"
	"fmt"

	"github.com/b-/gomox-uf/internal/clientinstantiator"
	"github.com/b-/gomox-uf/internal/pveurl"
	"github.com/b-/gomox-uf/internal/resourcesgetter"
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
			Name:        "vmid",
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
			Name: "storage",
			Usage: "Target storage for full clone.",
		},
		&cli.StringFlag{
			Name: "target",
			Usage: "Target node. Only allowed if the original VM is on shared storage.",
		},
	},
}

func bool2uint(b bool) uint {
	if b {
		return 1
	}
	return 0
}

// Clones a Proxmox VM as specified by the `from` arg and params struct
func cloneVm(c *cli.Context) error {
	client := clientinstantiator.InstantiateClient(
		pveurl.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	vmid := c.Uint64("vmid")
	vm, err := resourcesgetter.GetVirtualMachineByVMID(vmid, client, c.Context)
	if err != nil {
		return err
	}

	cloneOptions := proxmox.VirtualMachineCloneOptions {
		NewID: int(c.Uint64("newid")),
		BWLimit: c.Uint64("bwlimit"),
		Full: uint8(bool2uint(c.Bool("full"))),
		Name:     c.String("name"),
		Pool:     c.String("pool"),
		SnapName: c.String("snapname"),
		Storage:  c.String("storage"),
		Target:   c.String("target"),
	}

	if err != nil {
		return err
	}
	outvmid, task, err := vm.Clone(context.Background(), &cloneOptions)
	if err != nil {
		logrus.Panic("Oh no! ", err)
	}
	logrus.Info(fmt.Sprintf("clone requested! new id: %d.\n%#v\n", outvmid, task))
	return nil
}
