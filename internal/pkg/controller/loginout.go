package controller

import (
	"errors"
	"example/auth/internal/pkg/authorize"
	"example/auth/internal/pkg/models"
	"net/http"

	"example/auth/internal/pkg/database"
	"example/auth/internal/pkg/hashing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// function to handle home route
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hi mom : )",
	})
}

// function to handle login route
func Login(c *gin.Context) {
	// annonymous incoming req user
	var loginUser models.LoginUser
	c.ShouldBind(&loginUser)

	username := loginUser.Username
	password := loginUser.Password

	// actual user that is retrieved from the database
	var user models.User
	result := database.Db.Where("username=?", username).Find(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user doesn't exist"})
			return
		} else {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}
	}

	// else user exist
	ispwMatched := hashing.ComparePassword(password, user.Password)

	if !ispwMatched {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "password didn't match",
		})
		return
	}

	SignedAccessToken, saerr := authorize.GenerateAccessToken(models.User(user))

	if saerr != nil {
		c.AbortWithStatusJSON(498, gin.H{
			"message": "error signing jwt token",
		})
		return
	}

	SignedRefreshToken, expireRefresh, srerr := authorize.GenerateRefreshToken(models.User(user))

	if srerr != nil {
		c.AbortWithStatusJSON(498, gin.H{
			"message": "error signing jwt token",
		})
		return
	}

	// store the refreshToken in the database
	var isStored = authorize.StoreRefreshToken(database.Db, username, SignedRefreshToken, expireRefresh)

	if !isStored {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error storing the refresh token into the database",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":       "user logged in succesfully",
		"access token":  SignedAccessToken,
		"refresh token": SignedRefreshToken,
	})

}

// function to handle logout route
func Logout(c *gin.Context) {
	userAny, ok := c.Get("user")

	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting user from gin context"})
		return
	}

	user := userAny.(map[string]interface{})

	username := user["username"]

	var jwtRefreshTokenObj models.JwtRefreshToken
	// delete the refresh token
	resultrt := database.Db.Where("username=?", username).Delete(&jwtRefreshTokenObj)

	if resultrt.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error deleting refresh token"})
		return
	}

	// extract the access token
	jwtAccessToken, oktoken := c.Get("access_token")

	if !oktoken {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting jwt access token from gin.context"})
		return
	}

	// Add access token to blacklist
	at := models.JwtBlackListedToken{
		AccessToken: jwtAccessToken.(string),
	}

	resultat := database.Db.Create(&at)

	if resultat.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error inserting black list token"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "user logged out succesfully", "access token": jwtAccessToken})

}
