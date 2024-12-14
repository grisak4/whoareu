package login

import (
	"log"
	postgresqlmodels "whoareu/models/postgresql_models"
	"whoareu/utils/gen/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context, db *gorm.DB) {
	var loginform postgresqlmodels.User

	if err := c.BindJSON(&loginform); err != nil {
		c.JSON(400, gin.H{
			"error":   "incorrect data",
			"message": "error",
		})
		return
	}

	var user postgresqlmodels.User
	result := db.Where("email = ? AND password = ?", loginform.Email, loginform.Password).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		log.Printf("[ERROR] %s", result.Error)
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	jwtToken, err := jwt.GenerateJWT(user.ID, user.Username, user.RoleSystem)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	c.JSON(200, gin.H{
		"token": jwtToken,
	})
}
