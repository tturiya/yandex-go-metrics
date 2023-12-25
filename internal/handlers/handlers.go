package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/tturiya/iter5/internal/storage/memstorage"
)

var metrics = memstorage.NewMemStorage()

func UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {
	var (
		urlName  = chi.URLParam(r, "name")
		urlValue = chi.URLParam(r, "value")
	)

	urlValueInt64, err := strconv.ParseInt(urlValue, 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metrics.AddCounter(urlName, urlValueInt64)

	w.WriteHeader(http.StatusOK)
}

func UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		urlName  = chi.URLParam(r, "name")
		urlValue = chi.URLParam(r, "value")
	)

	urlValueFloat64, err := strconv.ParseFloat(urlValue, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metrics.AddGauge(urlName, urlValueFloat64)

	w.WriteHeader(http.StatusOK)
}

func GetCounterMetricHandler(w http.ResponseWriter, r *http.Request) {
	urlName := chi.URLParam(r, "name")
	if val, ok := metrics.GetCounter(urlName); ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", val)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetGaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	urlName := chi.URLParam(r, "name")
	if val, ok := metrics.GetGauge(urlName); ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%g", val)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseGlob("internal/views/home.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to parse HTML", err)
		return
	}

	data := struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}{
		Gauges:   metrics.GetAllGauges(),
		Counters: metrics.GetAllCounters(),
	}

	err = templ.ExecuteTemplate(w, "home", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to render HTML ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
