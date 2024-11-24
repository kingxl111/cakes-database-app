package main

import (
	"context"
	"log"

	"github.com/kingxl111/cakes-database-app/internal/storage/s3"

	"github.com/kingxl111/cakes-database-app/internal/logging"

	"time"

	"github.com/kingxl111/cakes-database-app/internal/config"
	"github.com/kingxl111/cakes-database-app/internal/server"
	"github.com/kingxl111/cakes-database-app/internal/service"
	"github.com/kingxl111/cakes-database-app/internal/storage"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//const (
//	envLocal = "local"
//	envDev   = "dev"
//	envProd  = "prod"
//)

func main() {
	ctx := context.Background()

	// config - cleanenv
	cfg := config.MustLoad()

	// logger
	//logg := SetupLogger(cfg.Env)
	logg, err := logging.NewLogger("logs.txt")
	if err != nil {
		log.Fatalf("can't initialize logger: %v", err.Error())
	}

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

	s3cl, err := s3.NewS3Client(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.Bucket,
		cfg.S3.Region,
		cfg.S3.PublicUrl)
	if err != nil {
		log.Fatalf("failed to connect to s3 client: %s", err)
	}

	// all layers
	st := storage.NewStorage(db)
	services := service.NewService(st, s3cl)
	router := server.NewHandler(services)

	// run server
	srv := &server.Server{}
	log.Printf("server started on %s", cfg.HTTPServer.Address)
	err = srv.Run(router.NewRouter(&ctx, logg.Lg, cfg.Env), cfg)
	if err != nil {
		logg.Lg.Info("server starting error!")
		return
	}
}

//func SetupLogger(env string) *slog.Logger {
//	var lg *slog.Logger
//
//	switch env {
//	case envLocal:
//		lg = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//		)
//	case envDev:
//		lg = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//		)
//	case envProd:
//		lg = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
//		)
//	default: // If env config is invalid, set prod settings by default due to security
//		lg = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
//		)
//	}
//
//	return lg
//}
