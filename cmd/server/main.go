package main

import (
	"fmt"
	"net/http"
	"server/internal/msg"
	"server/internal/server"
)

const ServerPort = "9479"

func init() {
	fmt.Println(msg.Logo)
}

func main() {
	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	middlewareMux := server.SetupMiddleware(mux)
	server.StartServer(middlewareMux, ServerPort)
}
