package list

import (
	"fmt"

	"github.com/b-/gomox/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/luthermonson/go-proxmox"
	"github.com/urfave/cli/v2"
)

const (
	Byte     = 1
	Kilobyte = 1024 * Byte
	Megabyte = 1024 * Kilobyte
	Gigabyte = 1024 * Megabyte
)

var Command = &cli.Command{
	Name:   "list",
	Usage:  "Lists virtual machines",
	Action: list,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "type",
			Category: "",
			// DefaultText: "both",
			FilePath:    "",
			Usage:       "`qemu|lxc|both`",
			Required:    false,
			Hidden:      false,
			HasBeenSet:  false,
			Value:       "both",
			Destination: nil,
			Aliases:     nil,
			EnvVars:     nil,
			TakesFile:   false,
			Action:      nil,
		},
	},
}

func list(c *cli.Context) error {
	client := util.InstantiateClient(
		util.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	rsList, err := util.GetVirtualMachineList(c.Context, client)
	if err != nil {
		return err
	}
	// simple table with zero customizations
	tw := table.NewWriter()
	// append a header row
	tw.AppendHeader(table.Row{"VMID", "Name", "Type", "Status", "Mem (MB)", "BootDisk (GB)", "PID"})
	// append some data rows

	for _, vm := range rsList {
		// if vm.

		tw.AppendRow(
			table.Row{
				// vmid,name,status,mem,boot,pid
				// https://git.proxmox.com/?p=qemu-server.git;a=blob;f=PVE/CLI/qm.pm;h=b17b4fe25d5bd21e9fe188e82998972b1dc29c36;hb=HEAD#l1001
				int(vm.VMID), vm.Name, vm.Status, vm.MaxMem / Megabyte,
				float64(vm.MaxDisk) / float64(Gigabyte),
				// uint64(vm.PID),
			},
		)
	}
	fmt.Print(tw.Render())
	return nil
}
