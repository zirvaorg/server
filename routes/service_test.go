package routes

import (
	"net/http"
	"server/internal/logic"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	mux := http.NewServeMux()
	Service(mux)

	tests := []struct {
		op    string
		param string
	}{
		{op: "ping", param: "zirva.org"},
		{op: "http", param: "zirva.org"},
		{op: "tcp", param: "zirva.org"},
		{op: "udp", param: "zirva.org"},
		{op: "traceroute", param: "zirva.org"},
	}

	for _, test := range tests {
		t.Run(test.op, func(t *testing.T) {
			resultChan := make(chan *logic.Response)
			go handleServiceOperation(test.op, test.param, resultChan)

			select {
			case res := <-resultChan:
				if res.Error != nil {
					t.Errorf("Error message: %s", res.Error.Error())
				}
				t.Logf("Response for %s: %+v", test.op, res)
				switch test.op {
				case "ping":
					if res.PingResult == nil {
						t.Errorf("expected PingResult, got nil")
					}

					t.Logf("PingResult: %+v", res.PingResult)
				case "http":
					if res.HttpResult == nil {
						t.Errorf("expected HttpResult, got nil")
					}

					t.Logf("HttpResult: %+v", res.HttpResult)
				case "tcp", "udp":
					if res.ConnectionResult == nil {
						t.Errorf("expected ConnectionResult, got nil")
					}

					t.Logf("ConnectionResult: %+v", res.ConnectionResult)
				case "traceroute":
					if res.TracerouteResult == nil {
						t.Errorf("expected TracerouteResult, got nil")
					}

					t.Logf("TracerouteResult: %+v", res.TracerouteResult)
				}

			case <-time.After(120 * time.Second):
				t.Errorf("operation %s timed out", test.op)
			}
		})
	}
}
