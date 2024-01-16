package database

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/tturiya/iter5/internal/handlers"
	"github.com/tturiya/iter5/internal/mertic"
)

type Database struct {
	file     *os.File
	interval int
}

func NewDatabase(inerval int, filename string) (*Database, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return &Database{
		interval: inerval,
		file:     file,
	}, nil
}

func (db *Database) StartLoop() {
	ticker := time.NewTicker(time.Duration(db.interval) *
		time.Second)
	for {
		select {
		case <-ticker.C:
			db.writeData()
		}
	}
}

// writes directly to metrics handler
func (db *Database) Consult() {
	var objs []*mertic.MetricsJSON
	data, err := io.ReadAll(db.file)

	err = json.Unmarshal(data, &objs)
	if err != nil {
		log.Println("Failed to byte->json, no writes to Metrics")
	}

	for _, x := range objs {
		switch x.MType {
		case "counter":
			handlers.Metrics.AddCounter(x.ID, *x.Delta)
		case "gauge":
			handlers.Metrics.AddGauge(x.ID, *x.Value)
		default:
			log.Fatalln("Unreachable code. Got metric type:", x.MType)
		}
	}
}

func (db *Database) writeData() {
	var (
		gauges   = handlers.Metrics.GetAllGauges()
		counters = handlers.Metrics.GetAllCounters()
		objs     []*mertic.MetricsJSON
	)

	for k, v := range gauges {
		objs = append(objs, &mertic.MetricsJSON{
			ID:    k,
			Value: &v,
			MType: "gauge",
			Delta: nil,
		})
	}

	for k, v := range counters {
		objs = append(objs, &mertic.MetricsJSON{
			ID:    k,
			Delta: &v,
			MType: "counter",
			Value: nil,
		})
	}

	b, err := json.Marshal(objs)
	if err != nil {
		log.Fatalln("DB crash on json->byte")
	}

	n, err := db.file.Write(b)
	if err != nil {
		log.Fatalln("DB crash on write")
	}

	log.Println("DB wrote", n, "bytes")
}
