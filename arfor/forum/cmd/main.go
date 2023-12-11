package main

import (
	"log"
	"os"

	"foruum"
	"foruum/handler"
	"foruum/repository"
	"foruum/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	configDB := repository.NewConfigDB()
	db, err := repository.InitDB(configDB)
	if err != nil {
		errorLog.Printf("Error initializing database : %v", err)
	}
	if err := repository.CreateTable(*db); err != nil {
		errorLog.Printf("Error creating table: %v", err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(foruum.Server)
	if err := srv.Run("1111", handlers.InitRoutes()); err != nil {
		errorLog.Printf("Error starting server : %v", err)
	}
}
