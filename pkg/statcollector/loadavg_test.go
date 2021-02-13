package statcollector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLoadAVG(t *testing.T) {
	ProcPath = "testdata/"
	avg, err := ParseLoadAVG()
	require.NoError(t, err)
	require.Equal(t, 0.50, avg.Load1)
	require.Equal(t, 0.64, avg.Load5)
	require.Equal(t, 0.65, avg.Load15)
}
