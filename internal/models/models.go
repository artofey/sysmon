package models

// StatSnapshot - Слепок всех видов статистики
type StatSnapshot struct {
	Lavg *LoadAVG
	Lcpu *LoadCPU
}

// LoadCPU - Средняя загрузка процессора
type LoadCPU struct {
	User   uint64
	System uint64
	Idle   uint64
}

// LoadAVG - Средняя загрузка системы за 1, 5 и 15 секунд
type LoadAVG struct {
	Load1  float64
	Load5  float64
	Load15 float64
}
