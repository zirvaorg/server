package server

import (
	"github.com/zirvaorg/ratelimit/memstore"
	"net/http"
	"server/routes"
	"time"
)

func SetupRoutes(mux *http.ServeMux) {
	registrarRateLimit := memstore.New(memstore.Options{
		Rate:            10 * time.Second,
		Limit:           3,
		BlockTime:       30 * time.Minute,
		CleanupInterval: 30 * time.Minute,
	})

	routes.Registrar(mux, registrarRateLimit)
	routes.Service(mux)
}
