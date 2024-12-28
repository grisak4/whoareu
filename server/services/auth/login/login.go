package login

import (
	"log"
	postgresqlmodels "whoareu/models/postgresql_models"
	"whoareu/utils/gen/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context, db *gorm.DB) {
	var loginform struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Парсинг JSON из запроса
	if err := c.BindJSON(&loginform); err != nil {
		c.JSON(400, gin.H{
			"error":   "incorrect data",
			"message": "error",
		})
		return
	}

	var user postgresqlmodels.User
	// Поиск пользователя по email и паролю в JSONB-поле
	result := db.Where("user_conf -> 'account' ->> 'email' = ? AND user_conf -> 'account' ->> 'password' = ?", loginform.Email, loginform.Password).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		log.Printf("[ERROR] %s", result.Error)
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	// Генерация JWT токена
	jwtToken, err := jwt.GenerateJWT(user.ID, user.UserConf, user.RoleSystem)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	// Возвращение токена
	c.JSON(200, gin.H{
		"token": jwtToken,
	})
}
