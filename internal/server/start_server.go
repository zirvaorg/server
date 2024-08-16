package server

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/logic"
	"server/internal/msg"
	"time"
)

func StartServer(handler http.Handler, port string) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logic.Output("info", fmt.Sprintf(msg.ServerRunning, port))

	if !logic.CheckAuthFile() {
		token := logic.GenerateRegistrarToken()
		logic.Output("warn", fmt.Sprintf(msg.RegistrarErr, port, token))
	}

	log.Fatal(server.ListenAndServe())
}
