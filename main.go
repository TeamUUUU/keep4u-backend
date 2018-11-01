package main

import (
	"context"
	"github.com/TeamUUUU/keep4u-backend/controllers"
	"github.com/TeamUUUU/keep4u-backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.NewClient("mongodb://mongo:27017")
	if err != nil {
		log.Fatal(err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := client.Connect(timeout); err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(timeout, nil);
		err != nil {
		log.Fatal(err)

	}
	boardsDao := services.BoardsDao{
		Db:             client,
		CollectionName: "boards",
		Database:       "keep4u-backend",
		Logger:         logger,
	}
	notesDAO := services.NotesDAO{
		Db:             client,
		CollectionName: "notes",
		Database:       "keep4u-backend",
		Logger:         logger,
	}
	api := controllers.ApiService{
		BoardsDAO: &boardsDao,
		NotesDAO:  &notesDAO,
		Logger:    logger,
	}
	r := gin.Default()
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://188.246.233.13:8080", "https://188.246.233.13:8080", "http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/boards", api.CreateBoard)
	r.POST("/boards/:board_id/notes", api.CreateNote)
	r.GET("/boards", api.GetUserBoards)
	r.GET("/boards/:board_id", api.GetBoard)
	r.GET("/boards/:board_id/notes", api.GetNotesForBoard)
	r.Run(":8080")
}
