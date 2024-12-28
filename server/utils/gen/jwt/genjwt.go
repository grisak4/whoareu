package jwt

import (
	"time"
	"whoareu/config/confget/jwtsec"
	postgresqlmodels "whoareu/models/postgresql_models"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId   uint `json:"user_id"`
	UserConf postgresqlmodels.UserConfig
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(user_id uint, userConf postgresqlmodels.UserConfig, role string) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		UserId:   user_id,
		UserConf: userConf,
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
