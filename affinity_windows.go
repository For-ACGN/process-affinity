package affinity

import (
	"golang.org/x/sys/windows"
)

// SetProcessAffinityByName is used to set process affinity by
// process name, it will set all process with the same name.
func SetProcessAffinityByName(name string, affinity uint32) error {
	processes, err := GetProcessList()
	if err != nil {
		return err
	}
	for i := 0; i < len(processes); i++ {
		if processes[i].Name == name {
			err = SetProcessAffinityByPID(processes[i].PID, affinity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetProcessAffinityByPID is used to set process affinity by PID.
func SetProcessAffinityByPID(pid uint32, affinity uint32) error {
	hProcess, err := OpenProcess(windows.PROCESS_SET_INFORMATION, false, pid)
	if err != nil {
		return err
	}
	defer CloseHandle(hProcess)

	return SetProcessAffinityMask(hProcess, affinity)
}
