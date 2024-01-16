package mertic

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/tturiya/iter5/internal/storage/memstorage"
)

type MetricsJSON struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func WriteMetric(ms memstorage.MetricsStorer) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	ms.AddGauge("MCacheSys", float64(memStats.MCacheSys))
	ms.AddGauge("Alloc", float64(memStats.Alloc))
	ms.AddGauge("BuckHashSys", float64(memStats.BuckHashSys))
	ms.AddGauge("Frees", float64(memStats.Frees))
	ms.AddGauge("GCCPUFraction", float64(memStats.GCCPUFraction))
	ms.AddGauge("GCSys", float64(memStats.GCSys))
	ms.AddGauge("HeapAlloc", float64(memStats.HeapAlloc))
	ms.AddGauge("HeapIdle", float64(memStats.HeapIdle))
	ms.AddGauge("HeapInuse", float64(memStats.HeapInuse))
	ms.AddGauge("HeapObjects", float64(memStats.HeapObjects))
	ms.AddGauge("HeapReleased", float64(memStats.HeapReleased))
	ms.AddGauge("HeapSys", float64(memStats.HeapSys))
	ms.AddGauge("LastGC", float64(memStats.LastGC))
	ms.AddGauge("Lookups", float64(memStats.Lookups))
	ms.AddGauge("MCacheInuse", float64(memStats.MCacheInuse))
	ms.AddGauge("MSpanSys", float64(memStats.MSpanSys))
	ms.AddGauge("MSpanInuse", float64(memStats.MSpanInuse))
	ms.AddGauge("Mallocs", float64(memStats.Mallocs))
	ms.AddGauge("NextGC", float64(memStats.NextGC))
	ms.AddGauge("NumForcedGC", float64(memStats.NumForcedGC))
	ms.AddGauge("NumGC", float64(memStats.NumGC))
	ms.AddGauge("OtherSys", float64(memStats.OtherSys))
	ms.AddGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
	ms.AddGauge("StackInuse", float64(memStats.StackInuse))
	ms.AddGauge("StackSys", float64(memStats.StackSys))
	ms.AddGauge("Sys", float64(memStats.Sys))
	ms.AddGauge("TotalAlloc", float64(memStats.TotalAlloc))
	ms.AddGauge("RandomValue", float64(rand.Intn(10)))
	ms.AddCounter("PollCount", int64(1))
}

func SendMetric(address string, ms memstorage.MetricsStorer) error {
	var (
		client = resty.New().
			SetHeader("Content-Type", "application/json")
		gauges   = ms.GetAllGauges()
		counters = ms.GetAllCounters()
		uri      = fmt.Sprintf("http://%s/update", address)
	)
	for key, val := range gauges {
		data := MetricsJSON{
			ID:    key,
			MType: "gauge",
			Value: &val,
		}
		_, err := client.R().SetBody(&data).Post(uri)
		if err != nil {
			return err
		}
	}

	for key, val := range counters {
		data := MetricsJSON{
			ID:    key,
			MType: "counter",
			Delta: &val,
		}
		_, err := client.R().SetBody(&data).Post(uri)
		if err != nil {
			return err
		}
	}
	return nil
}
