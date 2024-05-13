package main

import (
	"context"
	"fmt"
	"log"

	"github.com/godrrrain/exec-bash-api/src/config"
	"github.com/godrrrain/exec-bash-api/src/handler"
	"github.com/godrrrain/exec-bash-api/src/storage"
	"github.com/godrrrain/exec-bash-api/src/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("Error while reading config: %v. Set default values", err)
	}

	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Dbname, cfg.Database.Password)
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	aprocesses := types.NewActiveProcesses()

	handler := handler.NewHandler(psqlDB, aprocesses)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/commands/:uuid", handler.GetCommand)
	router.GET("/api/v1/commands", handler.GetCommands)
	router.POST("/api/v1/commands", handler.CreateCommand)
	router.POST("/api/v1/durables/commands", handler.CreateDurableCommand)
	router.PATCH("/api/v1/commands/:uuid", handler.StopCommand)
	router.DELETE("/api/v1/commands/:uuid", handler.DeleteCommand)

	router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
