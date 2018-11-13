package main

import (
	"context"
	"github.com/TeamUUUU/keep4u-backend/controllers"
	"github.com/TeamUUUU/keep4u-backend/middleware"
	"github.com/TeamUUUU/keep4u-backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/oauth2/v2"
	"log"
	"net/http"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.NewClient("mongodb://localhost:27017")
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
	accessService := services.DocumentAccessService{
		Db:             client,
		CollectionName: "access",
		Database:       "keep4u-backend",
		Logger:         logger,
	}
	api := controllers.ApiService{
		BoardsDAO:      &boardsDao,
		NotesDAO:       &notesDAO,
		Logger:         logger,
		DocumentAccess: &accessService,
	}

	oauthService, err := oauth2.New(&http.Client{})
	if err != nil {
		logger.Fatal("fail to create google oauth service", zap.Error(err))
	}
	authService := middleware.GoogleAuthMiddleware{
		Service: oauthService,
	}

	r := gin.Default()
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://188.246.233.13",  "https://keep4u.space", "http://keep4u.space","http://188.246.233.13", "http://188.246.233.13:8082", "https://188.246.233.13:8082", "http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(authService.Authorization())
	boards := r.Group("/boards")
	notes := r.Group("/notes")

	boards.POST("", api.CreateBoard)
	boardsId := boards.Group("/:board_id")
	boardsId.Use(middleware.Access(&accessService, "board_id"))
	boards.GET("", api.GetUserBoards)

	boardsId.GET("", api.GetBoard)
	boardsId.PATCH("", api.UpdateBoard)
	boardsId.DELETE("", api.DeleteBoard)

	boardsId.GET("/notes", api.GetNotesForBoard)
	boardsId.POST("/notes", api.CreateNote)

	boardsId.PATCH("/collaborators", api.UpdateBoardCollaborators)

	notesId := notes.Group("/:note_id")
	notesId.Use(middleware.Access(&accessService, "note_id"))
	notesId.DELETE("", api.DeleteNote)
	notesId.PATCH("", api.UpdateNote)
	notesId.GET("", api.GetNote)

	r.Run(":8080")
}
