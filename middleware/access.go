package middleware

import (
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/TeamUUUU/keep4u-backend/services"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
	"net/http"
)

func Access(accessService *services.DocumentAccessService, param string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, exists := ctx.Get("id_token")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Error{Message: "authorization required"})
			return
		}

		tokenInfo := token.(*oauth2.Tokeninfo)

		docId := ctx.Param(param)
		hasAccess, err := accessService.CheckAccess(&models.Access{UserID: tokenInfo.UserId, Documents: []string{docId}})
		if !hasAccess || err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.Error{Message: "resource unavailable"})
			return
		}
		accessService.Logger.Sugar().Debugw("got param", "err", err)
		ctx.Next()
	}
}
