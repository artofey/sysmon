// Сбор статистики и складывание ее в список

package statcollector

import (
	"fmt"
	"time"

	"github.com/artofey/sysmon/internal/pb"
)

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

// GetStat get snapsot all stats.
func GetStat() (*pb.StatSnapshot, error) {
	avg, err := ParseLoadAVG()
	if err != nil {
		return nil, fmt.Errorf("failed parse load average: %w", err)
	}
	cpu, err := ParseLoadCPU()
	if err != nil {
		return nil, fmt.Errorf("failed parse load process: %w", err)
	}
	s := pb.StatSnapshot{
		Lavg: avg,
		Lcpu: cpu,
	}
	return &s, nil
}

func StartColecting(out chan *pb.StatSnapshot) {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		s, err := GetStat()
		if err != nil {
			break
		}
		out <- s
	}
}
