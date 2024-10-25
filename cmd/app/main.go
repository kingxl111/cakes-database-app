package main

import (
	"cakes-database-app/pkg/config"
	server "cakes-database-app/pkg/http-server"
	"cakes-database-app/pkg/models"
	"cakes-database-app/pkg/storage/pgsql"
	"log"
)

func main() {
	// TODO: config - cleanenv
	cfg := config.MustLoad()
	log.Printf("%s", cfg.DB.Address)

	// TODO: database init
	db, err := pgsql.NewDB(
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

	u := models.User{
		FullName: "Medvedev Vadim Dmitrievich",
		Username: "kingxl111",
		Email: "kingxl111@mail.ru",
		PasswordHash: "askfdjhafkkasf",
		PhoneNumber: "+79xxxxxxx11",
	}	

	id, err := db.CreateUser(u)
	if err != nil {
		log.Fatalf("can't create user: %s", err.Error())
	}
	log.Println(id)

	// TODO: logger

	// TODO: router
	router := server.NewRouter()

	// TODO: middleware
	// router.Use() 

	// TODO: run server
	srv := &server.Server{}
	err = srv.Run(router, cfg)
	if err != nil {
		log.Fatal("server starting error!")
	}
}	