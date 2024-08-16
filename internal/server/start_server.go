package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal/logic"
	"server/internal/msg"
	"syscall"
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
		logic.Output("warn", fmt.Sprintf(msg.RegistrarErr, logic.ResolveExternalIP(), port, token))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logic.Output("error", fmt.Sprintf(err.Error()))
		}
	}()

	<-stop

	logic.Output("info", msg.ShutdownServer)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf(msg.ServerForceShutdown, err)
	}
}
