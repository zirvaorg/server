package routes

import (
	"net/http"
	"server/internal/logic"
	"server/middleware"
	"server/service"
	"time"
)

var serviceList = map[string]bool{
	"ping": true,
	"http": true, //@todo
	"tcp":  true, //@todo
}

func Service(mux *http.ServeMux) {
	mux.HandleFunc("GET /service/{op}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		op := r.PathValue("op")

		if !serviceList[op] {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch op {
		case "ping":
			p := r.URL.Query().Get("p")

			resolvedIP, err := logic.ResolveIP(p)
			if err != nil {
				logic.WriteResponse(w, &logic.Response{
					Status:  http.StatusBadRequest,
					Success: false,
				})
				return
			}

			ping, err := service.Ping(resolvedIP, 10, 10*time.Second)
			if err != nil {
				logic.WriteResponse(w, &logic.Response{
					Status:  http.StatusGatewayTimeout,
					Success: false,
				})
				return
			}

			logic.WriteResponse(w, &logic.Response{
				Status:     http.StatusOK,
				Success:    true,
				PingResult: &ping,
			})
		}
	}))
}
