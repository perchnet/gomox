package internal

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

type RequestableState string

const (
	RunningState = RequestableState(proxmox.StatusVirtualMachineRunning)
	StoppedState = RequestableState(proxmox.StatusVirtualMachineStopped)
	PausedState  = RequestableState(proxmox.StatusVirtualMachinePaused)
)

type StateRequestParams struct {
	RequestedState RequestableState
	Vm             *proxmox.VirtualMachine
}

// RequestState requests Proxmox change the state of a virtual machine.
func RequestState(params StateRequestParams, ctx context.Context) (*proxmox.Task, error) {

	var task *proxmox.Task
	var err error

	switch params.RequestedState {
	case RunningState:
		task, err = params.Vm.Start(ctx)
	case StoppedState:
		task, err = params.Vm.Stop(ctx)
	case PausedState:
		task, err = params.Vm.Pause(ctx)

	}
	logrus.Info(fmt.Sprintf("State %s requested! (vm: %d, task: %#v)", params.RequestedState, params.Vm.VMID, task))
	return task, err
}
