package controller

import (
	"example/auth/internal/pkg/database"
	"example/auth/internal/pkg/hashing"
	"example/auth/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// function to handle register route
func Register(c *gin.Context) {

	var user models.User
	binderr := c.Bind(&user)

	if binderr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// hash the password
	password, herr := hashing.Hash(user.Password)

	if herr != nil {
		panic("Error hashing the password")
	}

	// update password to hashed password and also UserId
	user.Password = password

	database.Db.Create(&user)

}

func RegisterOrganization(c *gin.Context) {
	var org models.Org
	binderr := c.Bind(&org)

	if binderr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	database.Db.Create(&org)
}
