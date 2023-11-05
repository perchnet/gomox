package staterequester

import (
	"context"
	"fmt"

	"github.com/b-/gomox-uf/internal/resourcesgetter"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)
type RequestableState string
const (
	RunningState = proxmox.StatusVirtualMachineRunning
	StoppedState = proxmox.StatusVirtualMachineStopped
	PausedState = proxmox.StatusVirtualMachinePaused
)
func RequestState(requestedState RequestableState,
	vmid uint64,
	client proxmox.Client,
	c context.Context,
	) (*proxmox.Task, error) {

	var task *proxmox.Task
	var err error

	vm,err := resourcesgetter.GetVirtualMachineByVMID(vmid, client, c)
	if err != nil {
		return nil, err
	}

	switch requestedState {
	case RunningState:
		task, err = vm.Start(context.Background())
	case StoppedState:
		task, err = vm.Stop(context.Background())
	case PausedState:
		task, err = vm.Pause(context.Background())

	}
	logrus.Info(fmt.Sprintf("State %s requested! (vm: %d, task: %#v)", requestedState, vm.VMID, task))
	return task,err
}
