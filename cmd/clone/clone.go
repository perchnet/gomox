package clone

import (
	"fmt"

	"github.com/b-/gomox/cmd/taskstatus"
	"github.com/b-/gomox/tasks"
	"github.com/b-/gomox/util"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

//goland:noinspection SpellCheckingInspection
var Command = &cli.Command{
	Name:   "clone",
	Usage:  "Clone a virtual machine.",
	Action: cloneVm,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:        "newid",
			Usage:       "`VMID` for the clone",
			Required:    false,
			Aliases:     []string{"n"},
			DefaultText: "next available",
		},

		&cli.Uint64Flag{
			Name:        "bwlimit",
			Usage:       "Override I/O bandwidth limit (in KiB/s).",
			DefaultText: "unlimited",
			Category:    "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "Description for the new VM.",
			Category: "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "format",
			Usage:    "Target format for file storage. Only valid for full clone. Can be raw, qcow, or vmdk.",
			Category: "Cloned VM Options:",
		},
		&cli.BoolFlag{
			Name:     "full",
			Usage:    "Create a full copy of all disks. This is always done when you clone a normal VM. For VM templates, we try to create a linked clone by default.",
			Category: "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "Set a name for the new VM.",
			Category: "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "pool",
			Usage:    "Add the new VM to the specified pool.",
			Category: "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "snapname",
			Usage:    "The name of the snapshot.",
			Category: "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "storage",
			Usage:    "Target storage for full clone.",
			Category: "Cloned VM Options:",
		},
		&cli.StringFlag{
			Name:     "target",
			Usage:    "Target node. Only allowed if the original VM is on shared storage.",
			Category: "Cloned VM Options:",
		},
		&cli.BoolFlag{
			Name:     "overwrite",
			Usage:    "Overwrite the target VMID if it already exists. (Note: only relevant when manually specifying VMID.)",
			Category: "Cloned VM Options:",
		},
	},
}

func bool2uint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// Clones a Proxmox VM as specified by the `from` arg and `params` struct
func cloneVm(c *cli.Context) error {
	newId := c.Uint64("newid")

	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)

	vmid, err := util.GetVmidArg(c.Args().Slice())
	if err != nil {
		return err
	}

	vm, err := util.GetVirtualMachineByVMID(c.Context, vmid, client)
	if err != nil {
		return err
	}
	if newId == 0 { // if newId isn't set
		newIdT, err := util.GetVmidArg(c.Args().Tail())
		if err == nil {
			newId = uint64(newIdT)
		}
	}
	cloneOptions := proxmox.VirtualMachineCloneOptions{
		NewID:    int(newId),
		BWLimit:  c.Uint64("bwlimit"),
		Full:     bool2uint8(c.Bool("full")),
		Name:     c.String("name"),
		Pool:     c.String("pool"),
		SnapName: c.String("snapname"),
		Storage:  c.String("storage"),
		Target:   c.String("target"),
	}

	if newId != 0 { // if we're manually assigning the target VMID
		vmWithSameId, _ := util.GetVirtualMachineByVMID(
			c.Context,
			newId,
			client,
		) // check if VM already exists with target VMID
		if vmWithSameId != nil {
			logrus.Infof("Virtual machine with target ID %d already exists.\n", newId)
			switch c.Bool("overwrite") {
			case true:
				task, err := util.DestroyVm(c.Context, vmWithSameId)
				if err != nil {
					return err
				}
				logrus.Info("overwrite requested\n")
				logrus.Warnf("destroying VM %d (%s)...\n", vmWithSameId.VMID, vmWithSameId.Name)
				logrus.Debugf("task: %s\n", task.UPID)

				// err = tasks.WaitTask(c.Context, task, tasks.WithSpinner())
				err = taskstatus.WaitForCliTask(c, &task)
				if err != nil {
					return err
				}

				logrus.Debugf("task: %s\n", task.UPID)
			case false:
				logrus.Tracef("%#v\n", vmWithSameId)
				return fmt.Errorf(
					"Use --overwrite if necessary.\n",
				)
			}

		}
	}

	newVmid, task, err := vm.Clone(c.Context, &cloneOptions) // do the clone
	if err != nil {
		return err
	}

	err = task.Ping(c.Context) // update task
	if err != nil {
		return err
	}
	if newVmid == 0 {
		newVmid = cloneOptions.NewID
	}

	logrus.Infof("clone requested! new id: %d.\n", newVmid)
	logrus.Tracef("%#v\n", task)
	if c.Bool("wait") {
		err = taskstatus.WaitForCliTask(c, task)
		if err != nil {
			return err
		}
	} else {
		logrus.Info(tasks.GetWaitCmd(*task))
	}

	return nil
}
