package logic

import (
	"encoding/json"
	"net/http"
	"server/service"
)

type Response struct {
	Status  int  `json:"status"`
	Success bool `json:"success"`

	PingResult *service.PingResult `json:"ping_result"`
}

func WriteResponse(w http.ResponseWriter, response *Response) {
	jsonBody, _ := json.Marshal(response)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(response.Status)
	_, err := w.Write(jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
