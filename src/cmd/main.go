package main

import (
	"context"
	"fmt"

	"github.com/godrrrain/exec-bash-api/src/handler"
	"github.com/godrrrain/exec-bash-api/src/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		"localhost", 5432, "postgres", "postgres", "postgres")
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	aprocesses := storage.NewActiveProcesses()

	handler := handler.NewHandler(psqlDB, aprocesses)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/commands/:uuid", handler.GetCommand)
	router.GET("/api/v1/commands", handler.GetCommands)
	router.POST("/api/v1/commands", handler.CreateCommand)

	router.Run(":8114")
}
