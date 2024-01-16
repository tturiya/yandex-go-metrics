package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tturiya/iter5/internal/config"
	"github.com/tturiya/iter5/internal/database"
	"github.com/tturiya/iter5/internal/handlers"
	"github.com/tturiya/iter5/internal/middleware/logger"
	"github.com/tturiya/iter5/internal/middleware/zipper"
)

func StartServer() error {
	var (
		r   = chi.NewRouter()
		cfg = config.NewServerConfig()
	)

	// set up middleware
	r.Use(logger.Logger)
	r.Use(middleware.Recoverer)
	r.Use(zipper.GZipper)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handlers.HomeHandler)

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetMetricsJSON)
		r.Get("/counter/{name}", handlers.GetCounterMetricHandler)
		r.Get("/gauge/{name}", handlers.GetGaugeMetricHandler)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.UpdateMetricsJSON)
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", handlers.BadRequestHandler)
	})

	db, err := database.NewDatabase(cfg.StoreInterval, cfg.StorageFP)
	if err != nil {
		log.Fatalln("Failed to initialize db:", err)
	}
	if cfg.Restore {
		db.Consult()
	}

	go db.StartLoop()

	return http.ListenAndServe(cfg.Addr, r)
}
