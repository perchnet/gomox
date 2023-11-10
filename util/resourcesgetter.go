package util

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
)

func GetVirtualMachineByVMID(ctx context.Context, vmid uint64, client proxmox.Client) (
	vm *proxmox.VirtualMachine,
	err error,
) {

	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx, "vm")
	if err != nil {
		return nil, err
	}

	for _, rs := range resources {
		if rs.VMID == vmid {
			node, err := client.Node(ctx, rs.Node)
			if err != nil {
				return nil, err
			}
			vm, err = node.VirtualMachine(ctx, int(rs.VMID))
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
