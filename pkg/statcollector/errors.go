package statcollector

import "fmt"

type ErrParseLoadCPU struct {
	Err error
}

func (e *ErrParseLoadCPU) Error() string {
	return fmt.Sprintf("ошибка при получении данных о загрузке процессора: %s", e.Err.Error())
}

type ErrParseLoadAVG struct {
	Err error
}

func (e *ErrParseLoadAVG) Error() string {
	return fmt.Sprintf("ошибка при получении данных о средней загрузке системы: %s", e.Err.Error())
}
