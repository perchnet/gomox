package util

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
)

const ( // resource getter filters
	VmFilter      = "vm"
	StorageFilter = "storage"
	NodeFilter    = "node"
	SdnFilter     = "sdn"
)

//goland:noinspection GoDeprecation
const ( // returned resource types to further filter
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

// getResourceListConfig defines the options for
type getResourceListConfig struct {
	filter        string
	furtherFilter []string
}

// GetResourceListOption specifies the type of Resources for GetResource to get.
type GetResourceListOption func(c *getResourceListConfig)

// WithVm makes GetResourceList return VMs.
func WithVm() GetResourceListOption {
	return func(c *getResourceListConfig) { c.filter = VmFilter }
}

// WithStorage makes GetResourceList return Storages.
func WithStorage() GetResourceListOption {
	return func(c *getResourceListConfig) { c.filter = StorageFilter }
}

// WithNode makes GetResourceList return Nodes.
func WithNode() GetResourceListOption {
	return func(c *getResourceListConfig) { c.filter = NodeFilter }
}

// WithSdn makes GetResourceList return SDNs.
func WithSdn() GetResourceListOption {
	return func(c *getResourceListConfig) { c.filter = SdnFilter }
}

// WithAll makes GetResourceList return VMs, Storage, Nodes, and SDNs.
func WithAll() GetResourceListOption {
	return func(c *getResourceListConfig) {
		c.filter = ""
	}
}

// WithQemu further filters GetResourceList for Qemu VMs.
func WithQemu() GetResourceListOption {
	return func(c *getResourceListConfig) { c.furtherFilter = append(c.furtherFilter, QemuResource) }
}

// WithLxc further filters GetResourceList for Qemu VMs.
func WithLxc() GetResourceListOption {
	return func(c *getResourceListConfig) { c.furtherFilter = append(c.furtherFilter, LxcResource) }
}

// WithPool further filters GetResourceList for Pools.
func WithPool() GetResourceListOption {
	return func(c *getResourceListConfig) { c.furtherFilter = append(c.furtherFilter, PoolResource) }
}

func GetResourceList(
	ctx context.Context,
	client proxmox.Client,
	opts ...GetResourceListOption,
) (
	rsList []*proxmox.ClusterResource,
	err error,
) {
	c := &getResourceListConfig{
		filter:        "",
		furtherFilter: nil,
	}
	for _, opt := range opts {
		opt(c)
	}
	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx, c.filter)
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
