package jwt

import (
	"time"
	"whoareu/config/confget/jwtsec"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(user_id uint, username, role string) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		UserId:   user_id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtsec.GetJwtToken())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
