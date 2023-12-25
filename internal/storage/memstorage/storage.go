package memstorage

type MetricsStorer interface {
	AddCounter(key string, counter int64)
	AddGauge(key string, gauge float64)
	GetCounter(key string) (int64, bool)
	GetGauge(key string) (float64, bool)
	GetAllCounters() map[string]int64
	GetAllGauges() map[string]float64
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms MemStorage) AddGauge(key string, gauge float64) {
	ms.gauge[key] = gauge
}

func (ms MemStorage) AddCounter(key string, counter int64) {
	if val, ok := ms.counter[key]; ok {
		ms.counter[key] = val + counter
	} else {
		ms.counter[key] = counter
	}
}

func (ms MemStorage) GetCounter(key string) (int64, bool) {
	val, ok := ms.counter[key]
	return val, ok
}

func (ms MemStorage) GetGauge(key string) (float64, bool) {
	val, ok := ms.gauge[key]
	return val, ok
}

func (ms MemStorage) GetAllCounters() map[string]int64 {
	return ms.counter
}

func (ms MemStorage) GetAllGauges() map[string]float64 {
	return ms.gauge
}
