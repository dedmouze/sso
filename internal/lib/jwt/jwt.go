package jwt

import (
	"fmt"
	"time"

	"sso/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	UID        int64
	Email      string
	Expiration time.Time
	AppID      int64
	Level      int8
}

// TODO: add tests
func NewToken(user models.User, app models.App, admin models.Admin, duration time.Duration) (string, error) {
	const op = "lib.jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID
	claims["level"] = admin.Level

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

// TODO: add tests
func Parse(raw, secret string) (*Token, error) {
	const op = "lib.jwt.Parse"

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token := &Token{
		UID:        int64(claims["uid"].(float64)),
		Email:      claims["email"].(string),
		Expiration: time.Unix(int64(claims["exp"].(float64)), 0),
		AppID:      int64(claims["app_id"].(float64)),
		Level:      int8(claims["level"].(float64)),
	}

	return token, nil
}
