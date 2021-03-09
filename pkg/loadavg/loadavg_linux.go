package loadavg

import (
	"fmt"
	"io/ioutil"

	"github.com/artofey/sysmon"
)

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse return load average info.
func (p *Parser) Parse() (interface{}, error) {
	procF := ProcPath + "loadavg"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, &ErrParseLoadAVG{fmt.Errorf("failed of read file %s: %w", procF, err)}
	}

	la := sysmon.LoadAVG{}
	fmt.Sscanf(string(b), "%f %f %f", &la.Load1, &la.Load5, &la.Load15)

	return &la, nil
}

// Average усредняет значения для массива значений LoadAVG.
func (Parser) Average(itemsI interface{}) (interface{}, error) {
	var l1, l5, l15 float64
	items, ok := itemsI.([]*sysmon.LoadAVG)
	if !ok {
		return nil, fmt.Errorf("itemsI type not is []*sysmon.LoadAVG")
	}
	for _, l := range items {
		l1 += l.Load1
		l5 += l.Load5
		l15 += l.Load15
	}
	count := float64(len(items))
	return &sysmon.LoadAVG{
		Load1:  l1 / count,
		Load5:  l5 / count,
		Load15: l15 / count,
	}, nil
}

func (Parser) Valid(item interface{}) bool {
	_, ok := item.(*sysmon.LoadAVG)
	return ok
}
