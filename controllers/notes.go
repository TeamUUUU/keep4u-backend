package controllers

import (
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *ApiService) CreateNote(ctx *gin.Context) {
	var noteCreate models.NoteCreate
	if err := ctx.BindJSON(&noteCreate);
		err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{Message: "unable to parse request"})
		return
	}
	noteCreate.BoardID = ctx.Param("board_id")
	note, err := api.NotesDAO.Create(&noteCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Error{Message: "unable to create note"})
		return
	}
	ctx.JSON(http.StatusCreated, note)
}

func (api *ApiService) GetNotesForBoard(ctx *gin.Context) {
	boardID := ctx.Param("board_id")
	if boardID == "" {
		ctx.JSON(http.StatusBadRequest, models.Error{Message: "board_id  parameter is missing"})
		return
	}
	boards, err := api.NotesDAO.GetNotesForBoard(boardID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Error{Message: "unable to find notes"})
		return
	}
	ctx.JSON(http.StatusOK, boards)
}
