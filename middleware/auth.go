package middleware

import (
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
	"net/http"
)

type GoogleAuthMiddleware struct {
	Service *oauth2.Service
}

func (gam *GoogleAuthMiddleware) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idToken := ctx.GetHeader("Authorization")
		resp, err := gam.Service.Tokeninfo().IdToken(idToken).Do()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Error{Message: "invalid authorization header"})
			return
		}
		ctx.Set("id_token", &resp)
		ctx.Next()
	}
}
