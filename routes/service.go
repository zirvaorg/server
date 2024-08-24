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
	"ping":       true,
	"http":       true,
	"tcp":        true,
	"udp":        true,
	"traceroute": true,
}

func handleServiceOperation(op, p string, resultChan chan<- *logic.Response) {
	defer close(resultChan)

	switch op {
	case "ping":
		resolvedIP, err := utils.ResolveIP(p)
		if err != nil {
			resultChan <- &logic.Response{
				Status:  http.StatusBadRequest,
				Success: false,
			}
			return
		}

		ping, err := service.Ping(resolvedIP, 10)
		if err != nil {
			resultChan <- &logic.Response{
				Status:  http.StatusGatewayTimeout,
				Success: false,
			}
			return
		}

		resultChan <- &logic.Response{
			Status:     http.StatusOK,
			Success:    true,
			PingResult: &ping,
		}

	case "http":
		httpResult, err := service.Http(p)
		resultChan <- &logic.Response{
			Status:     http.StatusOK,
			Success:    err == nil,
			HttpResult: &httpResult,
		}

	case "tcp", "udp":
		connectionResult, err := service.TcpOrUdp(p, op)
		resultChan <- &logic.Response{
			Status:           http.StatusOK,
			Success:          err == nil,
			ConnectionResult: &connectionResult,
		}

	case "traceroute":
		tracerouteResult, err := service.Traceroute(p)
		resultChan <- &logic.Response{
			Status:           http.StatusOK,
			Success:          err == nil,
			TracerouteResult: &tracerouteResult,
		}
	}
}

func Service(mux *http.ServeMux) {
	mux.HandleFunc("GET /service/{op}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		op := r.PathValue("op")

		if !serviceList[op] {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p := r.URL.Query().Get("p")
		resultChan := make(chan *logic.Response)

		go handleServiceOperation(op, p, resultChan)

		select {
		case res := <-resultChan:
			logic.WriteResponse(w, res)
		case <-time.After(120 * time.Second):
			logic.WriteResponse(w, &logic.Response{
				Status:  http.StatusGatewayTimeout,
				Success: false,
			})
		}
	}))
}
