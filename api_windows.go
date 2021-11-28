package affinity

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// reference:
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-setprocessaffinitymask

var (
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procSetProcessAffinityMask = modKernel32.NewProc("SetProcessAffinityMask")
)

// ProcessBasicInfo contains process basic information.
type ProcessBasicInfo struct {
	Name string
	PID  uint32
}

// GetProcessList is used to get process list that include PiD and name. // #nosec
func GetProcessList() ([]*ProcessBasicInfo, error) {
	const name = "GetProcessList"
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, newError(name, err, "failed to create process snapshot")
	}
	defer CloseHandle(snapshot)
	processes := make([]*ProcessBasicInfo, 0, 64)
	processEntry := &windows.ProcessEntry32{
		Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{})),
	}
	err = windows.Process32First(snapshot, processEntry)
	if err != nil {
		return nil, newError(name, err, "failed to call Process32First")
	}
	for {
		processes = append(processes, &ProcessBasicInfo{
			Name: windows.UTF16ToString(processEntry.ExeFile[:]),
			PID:  processEntry.ProcessID,
		})
		err = windows.Process32Next(snapshot, processEntry)
		if err != nil {
			if err.(windows.Errno) == windows.ERROR_NO_MORE_FILES {
				break
			}
			return nil, newError(name, err, "failed to call Process32Next")
		}
	}
	return processes, nil
}

// OpenProcess is used to open process by PID and return process handle.
func OpenProcess(desiredAccess uint32, inheritHandle bool, pid uint32) (windows.Handle, error) {
	const name = "OpenProcess"
	hProcess, err := windows.OpenProcess(desiredAccess, inheritHandle, pid)
	if err != nil {
		return 0, newErrorf(name, err, "failed to open process with PID %d", pid)
	}
	return hProcess, nil
}

// SetProcessAffinityMask is used to set process affinity mask with process handle.
func SetProcessAffinityMask(handle windows.Handle, mask uint32) error {
	const name = "SetProcessAffinityMask"
	ret, _, err := procSetProcessAffinityMask.Call(uintptr(handle), uintptr(mask))
	if ret == 0 {
		return newError(name, err, "failed to set process affinity")
	}
	return nil
}

// CloseHandle is used to close handle, it will not return error.
func CloseHandle(handle windows.Handle) {
	_ = windows.CloseHandle(handle)
}
