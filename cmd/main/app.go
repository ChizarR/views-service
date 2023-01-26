package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ChizarR/stats-service/internal/config"
	"github.com/ChizarR/stats-service/internal/intaraction"
	intrdb "github.com/ChizarR/stats-service/internal/intaraction/db"
	"github.com/ChizarR/stats-service/internal/user"
	userdb "github.com/ChizarR/stats-service/internal/user/db"
	"github.com/ChizarR/stats-service/pkg/client/mongodb"
	"github.com/ChizarR/stats-service/pkg/logging"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("CREATE router")
	mux := http.NewServeMux()

	cfg := config.GetConfig()
	cfgMongoDB := cfg.MongoDB

	logger.Info("CREATE Mongo Client")
	client, err := mongodb.NewClient(context.Background(), cfgMongoDB.Host, cfgMongoDB.Port,
		cfgMongoDB.User, cfgMongoDB.Password, cfgMongoDB.Database, cfgMongoDB.AuthDB)
	if err != nil {
		panic(err)
	}

	logger.Info("CREATE intaraction storage")
	intrStorage := intrdb.NewStorage(client, cfgMongoDB.Collections.Intaraction, logger)

	intrService := intaraction.NewIntaractionService(intrStorage, logger)

	logger.Info("REGISTER intaraction handler")
	intaractionHandler := intaraction.NewHandler(intrService, logger)
	intaractionHandler.Register(mux)

	logger.Info("CREATE user storage")
	userStorage := userdb.NewStorage(client, cfgMongoDB.Collections.User, logger)

	userService := user.NewUserStatService(userStorage, logger)

	logger.Info("REGISTER user handler")
	userHandler := user.NewHandler(userService, logger)
	userHandler.Register(mux)

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
