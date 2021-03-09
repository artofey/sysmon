package loadcpu

import "fmt"

type ErrParseLoadCPU struct {
	Err error
}

func (e *ErrParseLoadCPU) Error() string {
	return fmt.Sprintf("ошибка при получении данных о загрузке процессора: %s", e.Err.Error())
}
