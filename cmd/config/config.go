package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/b-/gomox/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "config",
	Usage:  "List the config settings of a Virtual Machine",
	Action: pveVersion,
	Flags:  []cli.Flag{},
}

func pveVersion(c *cli.Context) error {
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

	vm, err := util.GetVirtualMachineByVMID(c.Context, uint64(vmid), client)
	if err != nil {
		return err
	}

	// simple table with zero customizations
	tw := table.NewWriter()
	// append a header row
	tw.AppendHeader(table.Row{fmt.Sprintf("vm: %d", vmid), fmt.Sprintf("node: %s", vm.Node)})
	sets := make(map[string]*json.RawMessage)
	jThing, err := json.Marshal(vm.VirtualMachineConfig)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jThing, &sets)
	if err != nil {
		return err
	}
	for k, v := range sets {
		s := strings.Trim(string(*v), "\"")
		// append some data rows.
		tw.AppendRow(
			table.Row{
				k, s,
			},
		)

	}

	logrus.Infof("\n" + tw.Render())
	// fmt.Println(tw.Render())
	return nil
}
