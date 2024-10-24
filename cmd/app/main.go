package main

import (
	"cakes-database-app/pkg/config"
	server "cakes-database-app/pkg/http-server"
	"log"
)

func main() {
	// TODO: config - cleanenv
	cfg := config.MustLoad()
	log.Printf("%s", cfg.DB.Address)

	// TODO: database implementation!

	// TODO: logger

	// TODO: router
	router := server.NewRouter()

	// TODO: middleware
	// router.Use() 

	// TODO: run server
	srv := &server.Server{}
	err := srv.Run(router, cfg)
	if err != nil {
		log.Fatal("server starting error!")
	}
}