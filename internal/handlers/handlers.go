package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/tturiya/iter5/internal/mertic"
	"github.com/tturiya/iter5/internal/storage/memstorage"
	"github.com/tturiya/iter5/internal/util"
)

var Metrics = memstorage.NewMemStorage()

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
	Metrics.AddCounter(urlName, urlValueInt64)

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
	Metrics.AddGauge(urlName, urlValueFloat64)

	w.WriteHeader(http.StatusOK)
}

func GetCounterMetricHandler(w http.ResponseWriter, r *http.Request) {
	urlName := chi.URLParam(r, "name")
	if val, ok := Metrics.GetCounter(urlName); ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", val)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetGaugeMetricHandler(w http.ResponseWriter, r *http.Request) {
	urlName := chi.URLParam(r, "name")
	if val, ok := Metrics.GetGauge(urlName); ok {
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
		Gauges:   Metrics.GetAllGauges(),
		Counters: Metrics.GetAllCounters(),
	}

	err = templ.ExecuteTemplate(w, "home", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to render HTML ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var (
		contentType = r.Header["Content-Type"][0]
		accept      = "application/json"
		data        mertic.MetricsJSON
	)

	// include json type header in response
	w.Header().Set("Content-Type", accept)

	if contentType != accept {
		errStr := fmt.Sprintf("This path can't accept: %s", contentType)
		util.HTTPErrWrap(w, errStr, http.StatusNotFound, nil)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.HTTPErrWrap(w, "Failed to read request body",
			http.StatusNotFound, nil)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(b, &data)
	if err != nil {
		util.HTTPErrWrap(w, "Could not deserialize JSON on /update/",
			http.StatusNotFound, nil)
		return
	}

	switch strings.ToLower(data.MType) {
	case "gauge":
		if data.Value == nil {
			util.HTTPErrWrap(w, "Got nil gauge",
				http.StatusNotFound, nil)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		Metrics.AddGauge(data.ID, *data.Value)

	case "counter":
		if data.Delta == nil {
			util.HTTPErrWrap(w, "Got nil counter",
				http.StatusNotFound, nil)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		Metrics.AddCounter(data.ID, *data.Delta)

	default:
		util.HTTPErrWrap(w, "Unknown metric type",
			http.StatusNotFound, nil)
		return
	}
}

func GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var (
		contentType   = r.Header["Content-Type"][0]
		accept        = "application/json"
		clientAccepts = util.SliceContainsE[string](r.Header["Accept"],
			accept) || util.SliceContainsE[string](r.Header["Accept"],
			"*/*")
		data mertic.MetricsJSON
	)

	// include json type header in response
	w.Header().Add("Content-Type", accept)

	if r.ContentLength == 0 {
		util.HTTPErrWrap(w, "Got empty json",
			http.StatusNotFound, nil)
		return
	}

	if contentType != accept {
		errStr := fmt.Sprintf("This path can't accept: %s", contentType)
		util.HTTPErrWrap(w, errStr, http.StatusNotFound, nil)
		return
	}
	if !clientAccepts {
		util.HTTPErrWrap(w, "Client doesn't accept: ",
			http.StatusNotFound, &accept)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.HTTPErrWrap(w, "Failed to read request body",
			http.StatusNotFound, nil)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(b, &data)
	if err != nil {
		util.HTTPErrWrap(w, "Could not deserialize JSON on /value/",
			http.StatusNotFound, nil)
		return
	}

	switch strings.ToLower(data.MType) {
	case "gauge":
		v, ok := Metrics.GetGauge(data.ID)
		if !ok {
			util.HTTPErrWrap(w, fmt.Sprintf("Unknown gauge name: %s",
				data.ID),
				http.StatusNotFound, nil)
			return
		}

		respJSON := &mertic.MetricsJSON{
			MType: data.MType,
			ID:    data.ID,
			Value: &v,
			Delta: nil,
		}

		b, err := json.Marshal(respJSON)
		if err != nil {
			util.HTTPErrWrap(w, "Could not serialize JSON /value/gauge",
				http.StatusNotFound, nil)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			util.HTTPErrWrap(w, "Failed to write body",
				http.StatusInternalServerError, nil)
			return
		}

	case "counter":
		v, ok := Metrics.GetCounter(data.ID)
		if !ok {
			util.HTTPErrWrap(w, "Unknown gauge name",
				http.StatusNotFound, nil)
			return
		}

		respJSON := &mertic.MetricsJSON{
			MType: data.MType,
			ID:    data.ID,
			Delta: &v,
			Value: nil,
		}

		b, err := json.Marshal(respJSON)
		if err != nil {
			util.HTTPErrWrap(w, "Could not serialize JSON /value/counter",
				http.StatusNotFound, nil)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			util.HTTPErrWrap(w, "Failed to write body",
				http.StatusInternalServerError, nil)
			return
		}

	default:
		util.HTTPErrWrap(w, "Unknown metric type",
			http.StatusNotFound, nil)
		return
	}
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
