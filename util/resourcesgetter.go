package util

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
)

const (
	VmFilter      = "vm"
	StorageFilter = "storage"
	NodeFilter    = "node"
	SdnFilter     = "sdn"
)

//goland:noinspection GoDeprecation
const (
	NodeResource    = "node"
	StorageResource = "storage"
	PoolResource    = "pool"
	QemuResource    = "qemu"
	LxcResource     = "lxc"
	OpenVzResource  = "openvz" // deprecated
	SdnResource     = "sdn"
)

func GetVirtualMachineByVMID(ctx context.Context, vmid uint64, client proxmox.Client) (
	vm *proxmox.VirtualMachine,
	err error,
) {
	var node *proxmox.Node

	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx, VmFilter)
	if err != nil {
		return nil, err
	}

	for _, rs := range resources {
		if rs.VMID == vmid {
			node, err = client.Node(ctx, rs.Node)
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

func GetResourceList(ctx context.Context, client proxmox.Client, filter string) (
	rsList []*proxmox.ClusterResource,
	err error,
) {
	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx, filter)
	if err != nil {
		return nil, err
	}
	for _, rs := range resources {
		if err != nil {
			return nil, err
		}
		rsList = append(rsList, rs)
	}

	return rsList, nil
}

func GetVirtualMachineList(ctx context.Context, client proxmox.Client) (vmList []*proxmox.VirtualMachine, err error) {
	var node *proxmox.Node
	var vm *proxmox.VirtualMachine

	resources, err := GetResourceList(ctx, client, VmFilter)
	var rsList []*proxmox.ClusterResource
	for _, rs := range resources {
		node, err = client.Node(ctx, rs.Node)
		if err != nil {
			return nil, err
		}
		vm, err = node.VirtualMachine(ctx, int(rs.VMID))
		rsList = append(rsList, rs)
		if rs.Type == "qemu" {
			vmList = append(vmList, vm)
		}
	}

	return vmList, nil
}
