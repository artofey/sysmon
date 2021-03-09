package loadavg

import (
	"github.com/artofey/sysmon"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse return load average info.
func (p *Parser) Parse() (interface{}, error) {
	// todo
	return &sysmon.LoadAVG{}, nil
}

// AverageLoadAVG усредняет значения для массива значений LoadAVG.
func AverageLoadAVG(ll []*sysmon.LoadAVG) *sysmon.LoadAVG {
	// todo
	return nil
}
