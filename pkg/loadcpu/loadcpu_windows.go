package loadcpu

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/artofey/sysmon"
)

var cpuScanner *bufio.Scanner

func init() {
	cmd := exec.Command("typeperf", `\238(_Total)\*`)
	stdout, _ := cmd.StdoutPipe()

	cpuScanner = bufio.NewScanner(stdout)
	cpuScanner.Split(scanLastNonEmptyLine)
	err := cmd.Start()
	if err != nil {
		log.Fatal("Start error: ", err)
	}
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse return CPU stat.
func (p *Parser) Parse() (interface{}, error) {
	cpuScanner.Scan()
	record := cpuScanner.Text()
	st := strings.Split(record, ",")
	if len(st) < 10 || strings.HasPrefix(st[2], `"\\`) {
		return &sysmon.LoadCPU{}, nil
	}
	user, _ := strconv.ParseFloat(strings.Trim(st[2], `"`), 64)
	system, _ := strconv.ParseFloat(strings.Trim(st[3], `"`), 64)
	idle, _ := strconv.ParseFloat(strings.Trim(st[9], `"`), 64)
	if err := cpuScanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при сканировании вывода утилиты typeperf: %w", err)
	}
	return &sysmon.LoadCPU{
		User:   uint64(user),
		System: uint64(system),
		Idle:   uint64(idle),
	}, nil
}

// Average усредняет значения для массива значений LoadCPU.
func (Parser) Average(itemsI interface{}) (interface{}, error) {
	var lu, ls, li uint64
	items, ok := itemsI.([]*sysmon.LoadCPU)
	if !ok {
		return nil, fmt.Errorf("itemsI type not is []*sysmon..LoadCPU")
	}
	for _, l := range items {
		lu += l.User
		ls += l.System
		li += l.Idle
	}
	count := uint64(len(items))
	return &sysmon.LoadCPU{
		User:   lu / count,
		System: ls / count,
		Idle:   li / count,
	}
}

func (Parser) Valid(item interface{}) bool {
	_, ok := item.(*sysmon.LoadCPU)
	return ok
}

func scanLastNonEmptyLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		advance = len(data)
	} else {
		advance = bytes.LastIndexAny(data, "\n\r") + 1
	}
	data = data[:advance]

	data = bytes.TrimRight(data, "\n\r")

	if len(data) == 0 {
		return advance, nil, nil
	}

	token = data[bytes.LastIndexAny(data, "\n\r")+1:]
	return advance, token, nil
}
