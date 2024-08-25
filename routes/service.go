package routes

import (
	"net/http"
	"server/internal/logic"
	"server/middleware"
	"server/service"
	"sync"
	"time"
)

var serviceList = map[string]bool{
	"ping":       true,
	"http":       true,
	"tcp":        true,
	"udp":        true,
	"traceroute": true,
}

var mu sync.Mutex

func handleServiceOperation(op, p string, resultChan chan<- *logic.Response) {
	defer close(resultChan)

	switch op {
	case "ping":
		mu.Lock()
		defer mu.Unlock()

		ping, err := service.Ping(p, 10)
		resultChan <- &logic.Response{
			Success:    err == nil,
			Error:      err,
			PingResult: &ping,
		}

	case "http":
		httpResult, err := service.Http(p)
		resultChan <- &logic.Response{
			Success:    err == nil,
			Error:      err,
			HttpResult: &httpResult,
		}

	case "tcp", "udp":
		connectionResult, err := service.TcpOrUdp(p, op)
		resultChan <- &logic.Response{
			Success:          err == nil,
			Error:            err,
			ConnectionResult: &connectionResult,
		}

	case "traceroute":
		mu.Lock()
		defer mu.Unlock()

		tracerouteResult, err := service.Traceroute(p)
		resultChan <- &logic.Response{
			Success:          err == nil,
			Error:            err,
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
			res.Status = http.StatusOK

			if !res.Success {
				res.Status = http.StatusInternalServerError
				res.ErrorMessage = res.Error.Error()
				res.Error = nil
			}

			logic.WriteResponse(w, res)
		case <-time.After(120 * time.Second):
			logic.WriteResponse(w, &logic.Response{
				Status:  http.StatusGatewayTimeout,
				Success: false,
			})
		}
	}))
}
