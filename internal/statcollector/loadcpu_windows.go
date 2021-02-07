package statcollector

import (
	"github.com/artofey/sysmon/internal/pb"
)

// ParseLoadCPU return CPU stat.
func ParseLoadCPU() (*pb.LoadCPU, error) {
	// todo
	return nil, nil
}

// AverageLoadCPU усредняет значения для массива значений LoadCPU.
func AverageLoadCPU(ll []*pb.LoadCPU) *pb.LoadCPU {
	// todo
	return nil
}
