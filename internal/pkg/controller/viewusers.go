package controller

import (
	"example/auth/internal/pkg/database"
	"example/auth/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

/* Function to view users in a user's organization */
func ViewUsers(c *gin.Context) {
	userAny, ok := c.Get("user")

	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting user from gin context"})
		return
	}

	user := userAny.(map[string]interface{})
	// GET ORG ID OF THE USER
	orgid := user["orgid"]

	// DECLARE USERS ARRAY
	var users []models.User
	// EXTRACT ROWS WITH ORGID = orgid
	result := database.Db.Where("orgid=?", orgid).Find(&users)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error retrieving users from db with ORGID=orgid"})
		return
	}

	// extract the access token
	jwtAccessToken, oktoken := c.Get("access_token")

	if !oktoken {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting jwt access token from gin.context"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"users":        users,
		"access token": jwtAccessToken,
	})
}
