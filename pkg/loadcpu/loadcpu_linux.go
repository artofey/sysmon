package loadcpu

import (
	"fmt"
	"io/ioutil"

	"github.com/artofey/sysmon"
)

const multiplier = 1000000

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse return CPU stat.
func (p *Parser) Parse() (interface{}, error) {
	procF := ProcPath + "stat"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %w", procF, err)
	}

	lc := sysmon.LoadCPU{}
	var null float64
	fmt.Sscanf(string(b), "cpu %g %g %g %g", &lc.User, &null, &lc.System, &lc.Idle)
	lc.System /= multiplier
	lc.User /= multiplier
	lc.Idle /= multiplier

	return &lc, nil
}

// Average усредняет значения для массива значений LoadCPU.
func (Parser) Average(itemsI interface{}) (interface{}, error) {
	var lu, ls, li float64
	items, ok := itemsI.([]*sysmon.LoadCPU)
	if !ok {
		return nil, fmt.Errorf("itemsI type not is []*sysmon..LoadCPU")
	}
	for _, l := range items {
		lu += l.User
		ls += l.System
		li += l.Idle
	}
	count := float64(len(items))
	return &sysmon.LoadCPU{
		User:   lu / count,
		System: ls / count,
		Idle:   li / count,
	}, nil
}

func (Parser) Valid(item interface{}) bool {
	_, ok := item.(*sysmon.LoadCPU)
	return ok
}
