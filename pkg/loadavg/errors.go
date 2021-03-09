package loadavg

import "fmt"

type ErrParseLoadAVG struct {
	Err error
}

func (e ErrParseLoadAVG) Error() string {
	return fmt.Sprintf("ошибка при получении данных о средней загрузке системы: %s", e.Err.Error())
}
