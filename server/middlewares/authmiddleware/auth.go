package authmiddleware

import (
	"net/http"
	"whoareu/config/confget/jwtsec"
	jwtutils "whoareu/utils/gen/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Запрос без токена"})
			c.Abort()
			return
		}

		claims := &jwtutils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtsec.GetJwtToken(), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверная подпись токена"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка токена"})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			c.Abort()
			return
		}

		roleAllowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				roleAllowed = true
				break
			}
		}
		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Set("user_config", claims.UserConf)
		c.Next()
	}
}
