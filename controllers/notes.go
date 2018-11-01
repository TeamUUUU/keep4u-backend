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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "unable to parse request"})
		return
	}
	noteCreate.BoardID = ctx.Param("board_id")
	note, err := api.NotesDAO.Create(&noteCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "unable to create note"})
		return
	}
	ctx.JSON(http.StatusCreated, note)
}

func (api *ApiService) GetNotesForBoard(ctx *gin.Context) {
	boardID := ctx.Param("board_id")
	if boardID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "board_id  parameter is missing"})
		return
	}
	boards, err := api.NotesDAO.GetNotesForBoard(boardID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "unable to find notes"})
		return
	}
	ctx.JSON(http.StatusOK, boards)
}

func (api *ApiService) UpdateNote(ctx *gin.Context) {
	var noteUpdate models.NoteUpdate
	noteID := ctx.Param("note_id")
	if noteID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "note_id parameter is missing"})
		return
	}
	note, err := api.NotesDAO.UpdateNote(&noteUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{"fail to update note"})
		return
	}
	ctx.JSON(http.StatusOK, note)
}

func (api *ApiService) DeleteNote(ctx *gin.Context) {
	noteID := ctx.Param("note_id")
	if noteID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "note_id parameter is missing"})
		return
	}
	if err := api.NotesDAO.Delete(noteID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{"fail to delete note"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
