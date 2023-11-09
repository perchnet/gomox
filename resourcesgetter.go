package gomox

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
)

func GetVirtualMachineByVMID(vmid uint64, client proxmox.Client, c context.Context) (
	vm *proxmox.VirtualMachine,
	err error,
) {

	cluster, err := client.Cluster(c)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(c, "vm")
	if err != nil {
		return nil, err
	}

	for _, rs := range resources {
		if rs.VMID == vmid {
			node, err := client.Node(c, rs.Node)
			if err != nil {
				return nil, err
			}
			vm, err = node.VirtualMachine(c, int(rs.VMID))
			if err != nil {
				return nil, err
			}
		}
	}

	if vm == nil {
		err = fmt.Errorf("no vm with id found: %d", vmid)
	}

	return vm, err
}
