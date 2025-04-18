package token

import (
	"os"
	"strconv"
	"time"

	"github.com/agl/authentication/pkg/ip"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GetAccessToken(id string) (string, int, string, error) {
	godotenv.Load()
	
	PRIVATE_KEY := os.Getenv("PRIVATE_KEY")

	id_int, err := strconv.Atoi(id)

	if err != nil {
		return "", 0, "", err
	}

	tokenPairID := GetUniqueString()

	claims := &jwt.MapClaims{
		"iss": "Ilya",
		"exp": time.Now().Add(time.Hour).Unix(),
		"sub": id_int,
		"id":  tokenPairID,
		"data": map[string]string{
			"ip":    ip.GetLocalIP(),
			"email": "user@gmail.com",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	completeToken, err := token.SignedString([]byte(PRIVATE_KEY))
	if err != nil {
		return "", 0, "", err
	}

	return completeToken, id_int, tokenPairID, nil
}
