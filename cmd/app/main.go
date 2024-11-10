package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/kingxl111/cakes-database-app/internal/config"
	"github.com/kingxl111/cakes-database-app/internal/server"
	"github.com/kingxl111/cakes-database-app/internal/service"
	"github.com/kingxl111/cakes-database-app/internal/storage"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	// config - cleanenv
	cfg := config.MustLoad()

	// logger
	logg := SetupLogger(cfg.Env)

	// database init
	// wait 7 s before connect to db
	time.Sleep(7 * time.Second)
	db, err := storage.NewDB(
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLmode)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer db.Close()

	// all layers
	st := storage.NewStorage(db)
	services := service.NewService(st)
	router := server.NewHandler(services)

	// run server
	srv := &server.Server{}
	log.Printf("server started on %s", cfg.HTTPServer.Address)
	err = srv.Run(router.NewRouter(&ctx, logg, cfg.Env), cfg)
	if err != nil {
		logg.Info("server starting error!")
		return
	}
}

func SetupLogger(env string) *slog.Logger {
	var lg *slog.Logger

	switch env {
	case envLocal:
		lg = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		lg = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		lg = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		lg = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return lg
}
