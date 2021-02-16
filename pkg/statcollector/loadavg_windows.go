package statcollector

import (
	"github.com/artofey/sysmon"
)

// ParseLoadAVG return load average info.
func ParseLoadAVG() (*sysmon.LoadAVG, error) {
	// todo
	return &sysmon.LoadAVG{}, nil
}

// AverageLoadAVG усредняет значения для массива значений LoadAVG.
func AverageLoadAVG(ll []*sysmon.LoadAVG) *sysmon.LoadAVG {
	// todo
	return nil
}
