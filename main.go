package main

import (
	"github.com/TeamUUUU/keep4u-backend/controllers"
	"github.com/TeamUUUU/keep4u-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.NewClient("mongodb://localhost:27017")
	client.Connect(nil)
	if err != nil {
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
	}
	r := gin.Default()
	r.POST("/boards", api.CreateBoard)
	r.POST("/boards/:board_id/notes", api.CreateNote)
	r.GET("/boards", api.GetUserBoards)
	r.GET("/boards/:board_id", api.GetBoard)
	r.GET("/boards/:board_id/notes", api.GetNotesForBoard)
	r.Run()
}
