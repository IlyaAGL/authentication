package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/agl/authentication/internal/repositories"
	"github.com/agl/authentication/pkg/email"
	"github.com/agl/authentication/pkg/ip"
	"github.com/agl/authentication/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var savedRefreshToken string

func CreateTokens(id string, conn *pgx.Conn) (string, error) {
	accessToken, id_int, tokenPairID, err := token.GetAccessToken(id, "")

	refreshToken := token.GetUniqueString()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	err = repositories.StoreRefreshToken(id_int, string(hashedToken), tokenPairID, conn)

	if err != nil {
		return "", err
	}

	savedRefreshToken = refreshToken

	return accessToken, nil
}

func RefreshAccessToken(oldAccessToken string, conn *pgx.Conn) (string, error) {
	godotenv.Load()

	PRIVATE_KEY := os.Getenv("PRIVATE_KEY")

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(oldAccessToken, &claims, func(token *jwt.Token) (any, error) {
		return []byte(PRIVATE_KEY), nil
	})

	var userID int
	switch v := claims["sub"].(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		return "", errors.New("invalid user ID type")
	}

	current_ip := claims["data"].(map[string]any)["ip"].(string)
	current_email := claims["data"].(map[string]any)["email"].(string)
	current_tokenPairId := claims["id"].(string)

	receivedToken, actualTokenPairId := repositories.GetRefreshToken(userID, conn)

	if actualTokenPairId != current_tokenPairId {
		return "", errors.New("Not right token :(")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(receivedToken), []byte(savedRefreshToken)); err != nil {
		if current_ip != ip.GetLocalIP() {
			email.SendEmail(current_email)
		}

		return "", err
	} else {
		newAccessToken, _, _, err := token.GetAccessToken(strconv.Itoa(userID), actualTokenPairId)

		if err != nil {
			return "", errors.New("failed to create new token")
		}

		return newAccessToken, nil
	}
}
