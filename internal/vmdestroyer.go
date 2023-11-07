package internal

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

func DestroyVm(vm *proxmox.VirtualMachine, ctx context.Context) (*proxmox.Task, error) {
	task, err := vm.Delete(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Info(fmt.Sprintf("deletion requested! %#v", task))
	return task, nil
}

func DestroyVmWithForce(vm *proxmox.VirtualMachine, ctx context.Context) (*proxmox.Task, error) {
	/*
		task, err := vm.Stop(ctx)
		task, err := vm.Delete(ctx)
		if err != nil {
			return nil, err
		}
		logrus.Info(fmt.Sprintf("deletion requested! %#v", task))
		return task, nil
	*/
	return nil, fmt.Errorf(
		"Force deletion requested for VM %d.\n"+
			"Not implemented.", vm.VMID,
	)
}
