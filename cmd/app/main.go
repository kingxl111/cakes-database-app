package main

import (
	"cakes-database-app/pkg/config"
	server "cakes-database-app/pkg/server"
	"cakes-database-app/pkg/service"
	"cakes-database-app/pkg/storage"
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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
	db, err := storage.NewDB(
		cfg.DB.Username,
        cfg.DB.Password,
        cfg.DB.Address,
        cfg.DB.DBName,
        cfg.DB.SSLmode,
    )
    if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
    }
    defer db.Close()
	logg.Info("database started on", cfg.DB.Address)

	// err = Migrate(logg)	
	// if err != nil {
	// 	log.Fatalf("Migration up error: %s", err.Error())
	// }

	// all layers
	storage := storage.NewStorage(db)
	services := service.NewService(storage)
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

func Migrate(logg *slog.Logger) error {
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        logg.Error("Could not open database: ", err.Error())
		return err
    }

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        logg.Error("Could not create driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file:///home/vadim/cakes-database-app/pkg/storage/pgsql/migrations",
        "postgres", driver)
	if err != nil {
		logg.Error("Could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil {
		logg.Error("Could not apply migrations: %v", err)
	}

	logg.Info("Migrations applied successfully")

	version, dirty, err := m.Version()
	if err != nil {
		logg.Error("New migration version error: %s", err.Error())
	}

	logg.Info("Applied migration: %d, Dirty %t\n", version, dirty)
	return nil
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
