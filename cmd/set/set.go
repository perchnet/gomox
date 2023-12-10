package set

import (
	"strconv"
	"strings"

	"github.com/b-/gomox/cmd/taskstatus"
	"github.com/b-/gomox/util"
	"github.com/luthermonson/go-proxmox"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "set",
	Usage:  "Set virtual machine hardware",
	Action: set,
	Flags:  []cli.Flag{},
}

func set(c *cli.Context) error {
	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	vmid, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return err
	}

	vm, err := util.GetVirtualMachineByVMID(c.Context, uint64(vmid), client)
	if err != nil {
		return err
	}
	var options []proxmox.VirtualMachineOption
	for i := 0; i < len(c.Args().Tail()); i += 2 {
		k := c.Args().Tail()[i]
		k = strings.TrimPrefix(k, "-") // Allow two leading dashes
		k = strings.TrimPrefix(k, "-") // Allow two leading dashes

		v := c.Args().Tail()[i+1]
		v = strings.Trim(v, "\"") // remove quotes
		options = append(
			options, proxmox.VirtualMachineOption{
				Name:  k,
				Value: v,
			},
		)
	}
	task, err := vm.Config(c.Context, options...)
	if err != nil {
		return err
	}
	err = taskstatus.WaitForCliTask(c, task)
	if err != nil {
		return err
	}

	return nil
}
