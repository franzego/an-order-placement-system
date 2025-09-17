package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Auth struct {
	secretkey []byte
}

func NewAuthService(secret []byte) *Auth {
	return &Auth{secretkey: secret}
}

var jwtSecret = []byte("super-secret")

func (a *Auth) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	signedString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedString, nil

}

//var jwtSecret = []byte("super-secret")

func (a *Auth) VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, fmt.Errorf("invalid token")

}
