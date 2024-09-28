package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"server/internal/logic"
	"server/internal/msg"
	"server/internal/server"
)

const CurrentVersion = "0.0.3"

var DefaultServerPort = "9479"

func init() {
	v := flag.Bool("v", false, "show version")
	flag.StringVar(&DefaultServerPort, "p", DefaultServerPort, "server listen port")
	flag.Parse()

	if *v {
		fmt.Println(CurrentVersion)
		os.Exit(0)
	}

	fmt.Printf(msg.Logo, CurrentVersion)
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
