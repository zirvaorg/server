package main

import (
	"fmt"
	"net/http"
	"server/internal/logic"
	"server/internal/msg"
	"server/internal/server"
)

const ServerPort = "9479"

func init() {
	fmt.Println(msg.Logo)
}

func main() {
	if !logic.CheckPrivileges() {
		logic.Output("error", msg.PrivilegesErr)
		return
	}

	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	middlewareMux := server.SetupMiddleware(mux)
	server.StartServer(middlewareMux, ServerPort)
}
