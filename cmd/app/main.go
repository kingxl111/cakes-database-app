package main

import (
	"cakes-database-app/pkg/config"
	server "cakes-database-app/pkg/server"
	"cakes-database-app/pkg/service"
	"cakes-database-app/pkg/storage"
	"context"
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	ctx := context.Background()

	// TODO: config - cleanenv
	cfg := config.MustLoad()
	log.Printf("%s", cfg.DB.Address)

	// TODO: database init
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

	// err = Migrate()	
	// if err != nil {
	// 	log.Fatalf("Migration up error: %s", err.Error())
	// }

	// TODO: logger

	// all layers
	storage := storage.NewStorage(db)
	services := service.NewService(storage)
	router := server.NewHandler(services)
	
	// TODO: middleware
	// router.Use() 

	// TODO: run server
	srv := &server.Server{}
	err = srv.Run(router.NewRouter(&ctx), cfg)
	if err != nil {
		log.Fatal("server starting error!")
	}
}	

func Migrate() error {
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Could not open database: %v", err)
    }

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatalf("Could not create driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file:///home/vadim/cakes-database-app/pkg/storage/pgsql/migrations",
        "postgres", driver)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil {
		log.Fatalf("Could not apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("New migration version error: %s", err.Error())
	}

	log.Printf("Applied migration: %d, Dirty %t\n", version, dirty)
	return nil
}
