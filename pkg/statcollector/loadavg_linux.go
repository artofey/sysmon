package statcollector

import (
	"fmt"
	"io/ioutil"

	"github.com/artofey/sysmon"
)

// ParseLoadAVG return load average info.
func ParseLoadAVG() (*sysmon.LoadAVG, error) {
	procF := ProcPath + "loadavg"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %v", procF, err)
	}

	la := sysmon.LoadAVG{}
	fmt.Sscanf(string(b), "%f %f %f", &la.Load1, &la.Load5, &la.Load15)

	return &la, nil
}

// AverageLoadAVG усредняет значения для массива значений LoadAVG.
func AverageLoadAVG(ll []*sysmon.LoadAVG) *sysmon.LoadAVG {
	var l1, l5, l15 float64
	for _, l := range ll {
		l1 += l.Load1
		l5 += l.Load5
		l15 += l.Load15
	}
	return &sysmon.LoadAVG{
		Load1:  l1 / float64(len(ll)),
		Load5:  l5 / float64(len(ll)),
		Load15: l15 / float64(len(ll)),
	}
}
