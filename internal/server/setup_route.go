package server

import (
	"net/http"
	"server/routes"
)

func SetupRoutes(mux *http.ServeMux) {
	routes.Registrar(mux)
	routes.Service(mux)
}
