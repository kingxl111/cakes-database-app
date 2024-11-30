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

	// for authorizer database role
	authDB, err := storage.NewDB(
		cfg.DB.AuthUsername,
		cfg.DB.AuthPassword,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLmode)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	// explicit
	services.Authorization = service.NewAuthService(storage.NewStorage(authDB))

	// run server
	srv := &server.Server{}
	log.Printf("server started on %s", cfg.HTTPServer.Address)
	err = srv.Run(router.NewRouter(&ctx, logg.Lg, cfg.Env), cfg)
	if err != nil {
		logg.Lg.Info("server starting error!")
		return
	}
}
