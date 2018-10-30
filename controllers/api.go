package controllers

import "github.com/TeamUUUU/keep4u-backend/services"

type ApiService struct {
	BoardsDAO *services.BoardsDao
	NotesDAO  *services.NotesDAO
}
