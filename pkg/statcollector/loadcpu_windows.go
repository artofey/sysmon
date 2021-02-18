package statcollector

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

// ParseLoadCPU return CPU stat.
func ParseLoadCPU() (*sysmon.LoadCPU, error) {
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

// AverageLoadCPU усредняет значения для массива значений LoadCPU.
func AverageLoadCPU(ll []*sysmon.LoadCPU) *sysmon.LoadCPU {
	var lu, ls, li uint64
	for _, l := range ll {
		lu += l.User
		ls += l.System
		li += l.Idle
	}
	return &sysmon.LoadCPU{
		User:   lu / uint64(len(ll)),
		System: ls / uint64(len(ll)),
		Idle:   li / uint64(len(ll)),
	}
}

func scanLastNonEmptyLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Set advance to after our last line
	if atEOF {
		advance = len(data)
	} else {
		// data[advance:] now contains a possibly incomplete line
		advance = bytes.LastIndexAny(data, "\n\r") + 1
	}
	data = data[:advance]

	// Remove empty lines (strip EOL chars)
	data = bytes.TrimRight(data, "\n\r")

	// We have no non-empty lines, so advance but do not return a token.
	if len(data) == 0 {
		return advance, nil, nil
	}

	token = data[bytes.LastIndexAny(data, "\n\r")+1:]
	return advance, token, nil
}
