package config

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	vmid, err := strconv.Atoi(c.Args().First())
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
	// append some data rows.
	// config := *vm.VirtualMachineConfig
	/*v := reflect.ValueOf(config)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name
		val := v.Field(i).Interface()
		jVal, err := json.Marshal(val)
		if err != nil {
			return err
		}
		jName, err := json.Marshal(name)
		_ = jName
		if err != nil {
			return err
		}
		if len(jVal) > 2 && string(jVal) != "null" {

			// fmt.Printf("%s: %s\n", jName, val)
		}
	}

	*/
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
		tw.AppendRow(
			table.Row{
				k, s,
			},
		)

	}

	logrus.Infof(tw.Render())
	// fmt.Println(tw.Render())
	return nil
}
