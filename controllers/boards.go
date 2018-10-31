package controllers

import (
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *ApiService) CreateBoard(ctx *gin.Context) {
	var boardCreate models.BoardCreate
	if err := ctx.BindJSON(&boardCreate); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to parse request"})
		return
	}
	board, err := api.BoardsDAO.Create(&boardCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to create board"})
		return
	}
	ctx.JSON(http.StatusCreated, &board)
}

func (api *ApiService) GetUserBoards(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "user_id parameter missing"})
		return
	}
	boards, err := api.BoardsDAO.GetBoardsForUser(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to fetch user boards"})
		return
	}
	ctx.JSON(http.StatusOK, &boards)
}

func (api *ApiService) GetBoard(ctx *gin.Context) {
	boardID := ctx.Param("board_id")
	if boardID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "board_id parameter is missing"})
	}
	board, err := api.BoardsDAO.GetBoardById(boardID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to fetch board by id"})
		return
	}
	ctx.JSON(http.StatusOK, board)
}
