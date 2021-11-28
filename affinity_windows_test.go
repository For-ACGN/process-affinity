package affinity

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetProcessAffinityByName(t *testing.T) {
	notepad1 := exec.Command("notepad.exe")
	notepad2 := exec.Command("notepad.exe")

	err := notepad1.Start()
	require.NoError(t, err)
	err = notepad2.Start()
	require.NoError(t, err)

	defer func() {
		err = notepad1.Process.Kill()
		require.NoError(t, err)
		err = notepad2.Process.Kill()
		require.NoError(t, err)
	}()

	err = SetProcessAffinityByName("notepad.exe", 1)
	require.NoError(t, err)
}
