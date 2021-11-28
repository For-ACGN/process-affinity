package affinity

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func TestGetProcessList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		processes, err := GetProcessList()
		require.NoError(t, err)

		fmt.Println("Name    PID")
		for _, process := range processes {
			fmt.Printf("%s %d\n", process.Name, process.PID)
		}
	})
}

func TestOpenProcess(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hProcess, err := OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(os.Getpid()))
		require.NoError(t, err)

		CloseHandle(hProcess)
	})
}

func TestSetProcessAffinityMask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hProcess, err := OpenProcess(windows.PROCESS_SET_INFORMATION, false, uint32(os.Getpid()))
		require.NoError(t, err)
		defer CloseHandle(hProcess)

		err = SetProcessAffinityMask(hProcess, 1)
		require.NoError(t, err)
	})
}
