package controllers

import (
	"github.com/TeamUUUU/keep4u-backend/services"
	"go.uber.org/zap"
)

type ApiService struct {
	BoardsDAO      *services.BoardsDao
	NotesDAO       *services.NotesDAO
	DocumentAccess *services.DocumentAccessService
	Logger         *zap.Logger
}
