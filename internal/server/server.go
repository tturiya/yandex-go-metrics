package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tturiya/iter5/internal/config"
	"github.com/tturiya/iter5/internal/handlers"
)

func StartServer() error {
	var (
		r   = chi.NewRouter()
		cfg = config.NewServerConfig()
	)

	// set up middleware's
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handlers.HomeHandler)

	r.Route("/value", func(r chi.Router) {
		r.Get("/counter/{name}", handlers.GetCounterMetricHandler)
		r.Get("/gauge/{name}", handlers.GetGaugeMetricHandler)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", handlers.BadRequestHandler)
	})

	return http.ListenAndServe(cfg.Addr, r)
}
