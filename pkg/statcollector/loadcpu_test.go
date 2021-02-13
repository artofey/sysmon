package statcollector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLoadCPU(t *testing.T) {
	ProcPath = "testdata/"
	cpu, err := ParseLoadCPU()
	require.NoError(t, err)
	require.Equal(t, uint64(7037598), cpu.User)
	require.Equal(t, uint64(3377528), cpu.System)
	require.Equal(t, uint64(464549088), cpu.Idle)
}
