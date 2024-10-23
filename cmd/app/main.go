package main

import (
	"cakes-database-app/pkg/config"
	"log"
)



func main() {
	// TODO: config - cleanenv
	cfg := config.MustLoad()
	log.Printf("%s", cfg.DB.Address)

	// TODO: database implementation!

	// TODO: logger

	// TODO: router

	// TODO: run server
}