package api

import (
	"net/http"

	"github.com/agl/authentication/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(conn *pgx.Conn) {
	r := gin.Default()

	r.GET("/tokens/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		accessToken, err := services.CreateTokens(id, conn)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Something went wrong :(",
			})

			return
		}

		ctx.SetCookie(
			"access_token",
			accessToken,
			3600,
			"/",
			"",
			false,
			false,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"status": "Ok",
		})

		return
	})

	r.GET("/tokens", func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie("access_token")

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access token required",
			})

			return
		}

		newToken, err := services.RefreshAccessToken(accessToken, conn)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": err.Error(),
			})

			return
		}

		ctx.SetCookie(
			"access_token",
			newToken,
			3600,
			"/",
			"",
			false,
			false,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"status": "Ok",
		})

		return
	})

	r.Run(":6060")
}
