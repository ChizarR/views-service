package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ChizarR/stats-service/internal/config"
	"github.com/ChizarR/stats-service/internal/intaraction"
	"github.com/ChizarR/stats-service/internal/intaraction/db"
	"github.com/ChizarR/stats-service/pkg/client/mongodb"
	"github.com/ChizarR/stats-service/pkg/logging"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("CREATE router")
	mux := http.NewServeMux()

	cfg := config.GetConfig()
	cfgMongoDB := cfg.MongoDB

	logger.Info("CREATE Mongo Client...")
	client, err := mongodb.NewClient(context.Background(), cfgMongoDB.Host, cfgMongoDB.Port,
		cfgMongoDB.User, cfgMongoDB.Password, cfgMongoDB.Database, cfgMongoDB.AuthDB)
	if err != nil {
		panic(err)
	}

	logger.Info("CREATE storage...")
	storage := db.NewStorage(client, cfgMongoDB.Collection, logger)

	intrService := intaraction.NewIntaractionService(storage, logger)

	logger.Info("REGISTER intaraction handler")
	intaraction_handler := intaraction.NewHandler(intrService, logger)
	intaraction_handler.Register(mux)

	run(mux, cfg)
}

func run(mux *http.ServeMux, cfg *config.Config) {
	logger := logging.GetLogger()

	logger.Info("STARTING application")

	addr := fmt.Sprintf("%s:%s", cfg.Server.BindIP, cfg.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("SERVER is listening on http://%s", addr)
	logger.Fatal(server.Serve(listener))
}
