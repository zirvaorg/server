package routes

import (
	"net/http"
	"server/internal/logic"
	"server/internal/utils"
	"server/middleware"
	"server/service"
	"time"
)

var serviceList = map[string]bool{
	"ping": true,
	"http": true,
	"tcp":  true,
	"udp":  true,
}

func Service(mux *http.ServeMux) {
	mux.HandleFunc("GET /service/{op}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		op := r.PathValue("op")

		if !serviceList[op] {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p := r.URL.Query().Get("p")

		switch op {
		case "ping":
			resolvedIP, err := utils.ResolveIP(p)
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

		case "http":
			httpResult, err := service.Http(p)
			logic.WriteResponse(w, &logic.Response{
				Status:     http.StatusOK,
				Success:    err == nil,
				HttpResult: &httpResult,
			})

		case "tcp", "udp":
			connectionResult, err := service.TcpOrUdp(p, op)
			logic.WriteResponse(w, &logic.Response{
				Status:           http.StatusOK,
				Success:          err == nil,
				ConnectionResult: &connectionResult,
			})
		}
	}))
}
