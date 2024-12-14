package regin

import (
	"log"
	"net/http"
	"time"
	postgresqlmodels "whoareu/models/postgresql_models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Regin(c *gin.Context, db *gorm.DB) {
	var reginform postgresqlmodels.User

	if err := c.ShouldBindJSON(&reginform); err != nil {
		log.Printf("[ERROR] %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	reginform.RoleSystem = "user"
	reginform.Created_At = time.Now()
	reginform.Updated_At = time.Now()

	if err := db.Create(&reginform).Error; err != nil {
		log.Printf("[ERROR] %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}
