package util

import (
	"fmt"
	"strconv"
)

const (
	MaxVmid = 999999999
	MinVmid = 100
)

func VmidOutOfRangeError() error {
	return fmt.Errorf("please supply a VMID between %d and %d", MinVmid, MaxVmid)
}
func CheckVmidRange(vmid uint64) error {
	if vmid < MinVmid || vmid > MaxVmid {
		return VmidOutOfRangeError()
	}
	return nil
}

func GetVmidArg(args []string) (uint64, error) {
	if len(args) == 0 {
		return 0, VmidOutOfRangeError()
	}
	ivmid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, VmidOutOfRangeError()
	}
	err = CheckVmidRange(uint64(ivmid))
	if err != nil {
		return 0, VmidOutOfRangeError()
	}
	vmid := uint64(ivmid)
	return vmid, nil
}
