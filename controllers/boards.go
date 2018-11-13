package controllers

import (
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/oauth2/v2"
	"net/http"
)

func (api *ApiService) CreateBoard(ctx *gin.Context) {
	var boardCreate models.BoardCreate
	if err := ctx.BindJSON(&boardCreate); err != nil {
		api.Logger.Error("fail to bind json params", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to parse request"})
		return
	}
	ownerIDraw, exists := ctx.Get("id_token")
	if !exists {
		api.Logger.Error("user_id not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "user_id parameter missing"})
		return
	}
	ownerID := ownerIDraw.(*oauth2.Tokeninfo).UserId
	boardCreate.OwnerID = ownerID
	board, err := api.BoardsDAO.Create(&boardCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to create board"})
		return
	}
	if err := api.DocumentAccess.UpdateAccess(&models.Access{UserID: ownerID, Documents: []string{board.ID}});
		err != nil {
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
		return
	}
	board, err := api.BoardsDAO.GetBoardById(boardID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to fetch board by id"})
		return
	}
	ctx.JSON(http.StatusOK, board)
}

func (api *ApiService) UpdateBoard(ctx *gin.Context) {
	boardID := ctx.Param("board_id")
	if boardID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "board_id parameter is missing"})
		return
	}
	var boardUpdate models.BoardUpdate
	if err := ctx.BindJSON(&boardUpdate); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to parse request body"})
		return
	}
	boardUpdate.ID = boardID
	board, err := api.BoardsDAO.Update(&boardUpdate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.AbortWithStatusJSON(http.StatusNotFound, models.Error{Message: "board with such id not found"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to update board"})
		return
	}
	ctx.JSON(http.StatusOK, board)
}

func (api *ApiService) UpdateBoardCollaborators(ctx *gin.Context) {
	boardID := ctx.Param("board_id")
	if boardID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "board_id parameter is missing"})
		return
	}
	var update models.BoardCollaborationUpdate
	if err := ctx.BindJSON(&update); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to parse request body"})
		return
	}
	updatedCollaboration, err := api.BoardsDAO.AddCollaborators(boardID, update)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to update collaboration"})
		return
	}
	ctx.JSON(http.StatusOK, updatedCollaboration)
}

func (api *ApiService) DeleteBoard(ctx *gin.Context) {
	boardID := ctx.Param("board_id")
	if boardID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "board_id parameter is missing"})
		return
	}
	if err := api.BoardsDAO.Delete(boardID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to delete board"})
		return
	}
	if err := api.NotesDAO.DeleteByBoardId(boardID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to delete board notes"})
		return
	}
	ownerIDraw, exists := ctx.Get("id_token")
	if !exists {
		api.Logger.Error("user_id not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "user_id parameter missing"})
		return
	}
	ownerID := ownerIDraw.(*oauth2.Tokeninfo).UserId
	if err := api.DocumentAccess.DropAccess(&models.Access{UserID: ownerID, Documents: []string{boardID}});
		err != nil {
		api.Logger.Error("user_id not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to drop a board"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
