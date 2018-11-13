package controllers

import (
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/api/oauth2/v2"
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
	token, exists := ctx.Get("id_token")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "id_token expected"})
		return
	}
	tokenInfo := token.(*oauth2.Tokeninfo)
	noteCreate.OwnerID = tokenInfo.UserId
	note, err := api.NotesDAO.Create(&noteCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to create note"})
		return
	}
	if err := api.DocumentAccess.UpdateAccess(&models.Access{UserID: tokenInfo.UserId, Documents: []string{note.ID}});
		err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to create note"})
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

func (api *ApiService) GetNote(ctx *gin.Context) {
	noteID := ctx.Param("note_id")
	if noteID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "note_id parameter is missing"})
		return
	}
	note, err := api.NotesDAO.GetNote(noteID)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			ctx.AbortWithStatusJSON(http.StatusNotFound, models.Error{Message: "note_id not found"})
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: "fail to find note by note_id"})
			return
		}
	}
	ctx.JSON(http.StatusOK, note)
}

func (api *ApiService) UpdateNote(ctx *gin.Context) {
	var noteUpdate models.NoteUpdate
	noteID := ctx.Param("note_id")
	if noteID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "note_id parameter is missing"})
		return
	}
	if err := ctx.BindJSON(&noteUpdate); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to parse json"})
		return
	}

	noteUpdate.ID = noteID
	note, err := api.NotesDAO.Update(&noteUpdate)
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
		if err == mongo.ErrNoDocuments {
			ctx.AbortWithStatusJSON(http.StatusNotFound, models.Error{Message: "note with such id doesn't exists"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{"fail to delete note"})
		return
	}
	ownerIDraw, exists := ctx.Get("id_token")
	if !exists {
		api.Logger.Error("user_id not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "user_id parameter missing"})
		return
	}
	ownerID := ownerIDraw.(*oauth2.Tokeninfo).UserId
	if err := api.DocumentAccess.DropAccess(&models.Access{UserID: ownerID, Documents: []string{noteID}});
		err != nil {
		api.Logger.Error("user_id not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: "fail to drop a board"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
