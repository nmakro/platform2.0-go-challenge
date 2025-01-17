package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/nmakro/platform2.0-go-challenge/environment"
	"github.com/nmakro/platform2.0-go-challenge/internal/session"
	"github.com/nmakro/platform2.0-go-challenge/server"
	assetModule "github.com/nmakro/platform2.0-go-challenge/server/modules/assets"
	userModule "github.com/nmakro/platform2.0-go-challenge/server/modules/user"

	"github.com/spf13/viper"
)

var Version = "development"

func main() {
	environment.LoadConfig("")
	app := NewApp()
	router := mux.NewRouter()

	sessionStore := session.GetSessionStore()
	assetModule.Setup(router, app.assetService, sessionStore)
	userModule.Setup(router, app.userService, sessionStore)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf(Version)
	s := server.Server{
		Router: router,
		HttpServer: &http.Server{
			Addr:         viper.GetString("SERVER_ADDRESS"),
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     logger,
		},
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := s.HttpServer.Shutdown(ctx)
		if err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server running at:", viper.GetString("SERVER_ADDRESS"))

	if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server error: %s", err.Error())
	}

	<-done
	logger.Println("Server stopped")
}
