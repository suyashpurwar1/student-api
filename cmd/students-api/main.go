package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suyashpurwar1/students-api/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Setup router
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students API"))
	})
  
	// Setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}


	fmt.Println("Server started at:", cfg.HTTPServer.Addr)

	slog.Info(("server started"),slog.String("address",cfg.Addr))
	
	done := make(chan os.Signal,1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT ,syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	err:=server.Shutdown(ctx)
	if err!=nil{
		slog.Error("Failed to shutdown server",slog.String("error",err.Error()))
	} 

	slog.Info("server shutdown successfully")
}
