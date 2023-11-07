package internal

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

type DestroyParams struct {
	RequestedState RequestableState
	Vm             *proxmox.VirtualMachine
}

func DestroyVm(params DestroyParams, ctx context.Context) (*proxmox.Task, error) {
	task, err := params.Vm.Delete(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Info(fmt.Sprintf("deletion requested! %#v", task))
	return task, nil
}
