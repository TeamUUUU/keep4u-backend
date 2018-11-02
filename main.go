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
		logger.Fatal("fail to setup mongo client", zap.Error(err))
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := client.Connect(timeout); err != nil {
		logger.Fatal("fail to connet to mongo", zap.Error(err))
	}
	if err := client.Ping(timeout, nil);
		err != nil {
		logger.Fatal("fail to ping mongo", zap.Error(err))

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
		AllowOrigins:     []string{"http://188.246.233.13:8082", "https://188.246.233.13:8082", "http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	boards := r.Group("/boards")
	notes := r.Group("/notes")

	boards.POST("/", api.CreateBoard)
	boards.GET("/", api.GetUserBoards)

	boards.GET("/:board_id", api.GetBoard)
	boards.PUT("/:board_id", api.UpdateBoard)
	boards.DELETE("/:board_id", api.DeleteBoard)

	boards.GET("/:board_id/notes", api.GetNotesForBoard)
	boards.POST("/:board_id/notes", api.CreateNote)

	boards.PATCH("/:board_id/collaborators", api.UpdateBoardCollaborators)

	notes.DELETE("/:note_id", api.DeleteNote)
	notes.PATCH("/:note_id", api.UpdateNote)
	notes.GET("/:note_id", api.GetNote)

	r.Run(":8080")
}
