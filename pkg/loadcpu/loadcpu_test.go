package loadcpu

import (
	"testing"

	"github.com/artofey/sysmon"
	"github.com/stretchr/testify/require"
)

func TestParseLoadCPU(t *testing.T) {
	ProcPath = "testdata/"
	p := NewParser()
	cpuI, err := p.Parse()
	cpu := cpuI.(*sysmon.LoadCPU)
	require.NoError(t, err)
	require.Equal(t, float64(7037598)/multiplier, cpu.User)
	require.Equal(t, float64(3377528)/multiplier, cpu.System)
	require.Equal(t, float64(464549088)/multiplier, cpu.Idle)
}
