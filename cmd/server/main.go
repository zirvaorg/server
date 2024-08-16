package main

import (
	"flag"
	"fmt"
	"net/http"
	"server/internal/logic"
	"server/internal/msg"
	"server/internal/server"
)

var DefaultServerPort = "9479"

func init() {
	fmt.Println(msg.Logo)

	flag.StringVar(&DefaultServerPort, "p", DefaultServerPort, "server listen port")
	flag.Parse()
}

func main() {
	err := logic.CheckEnvironment(DefaultServerPort)
	if err != nil {
		logic.Output("error", err.Error())
		return
	}

	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	middlewareMux := server.SetupMiddleware(mux)
	server.StartServer(middlewareMux, DefaultServerPort)
}
