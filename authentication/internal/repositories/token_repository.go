package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func StoreRefreshToken(id int, hashedToken, tokenPairId string, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO tokens VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET token = $2, tokenPairID = $3",
		id, hashedToken, tokenPairId)

	return err
}

func GetRefreshToken(id int, conn *pgx.Conn) (string, string) {
	var recievedRefreshToken string
	var receivedTokenPairID string
	conn.QueryRow(context.Background(),
		"SELECT token, tokenPairID FROM tokens WHERE id=$1", id).Scan(&recievedRefreshToken, &receivedTokenPairID)

	return recievedRefreshToken, receivedTokenPairID
}
