package controller

import (
	"fmt"
	"net/http"

	"example/auth/internal/pkg/database"
	"example/auth/internal/pkg/hashing"
	"example/auth/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

/* [ ] Middleware function to check if admin user made the request */
func CheckIfAdminUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, ok := c.Get("user")

		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting user userValue"})
			return
		}

		user := userAny.(map[string]interface{})

		isAdmin := user["admin"].(bool)

		// if non admin user requested then abort
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"message": "non admin user can't add new user",
			})
			return
		}

		c.Set("orgid", user["orgid"])

		// call the handler
		c.Next()

	}
}

/* [ ] Function to add new user to an organization */
func AddUser(c *gin.Context) {
	orgId, ok := c.Get("orgid")

	if !ok {
		fmt.Println("orgId : ", orgId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting orgid"})
		return
	}

	orgid := uint(orgId.(float64))

	// extract the access token
	jwtAccessToken, oktoken := c.Get("access_token")

	if !oktoken {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting jwt access token from gin.context"})
		return
	}

	var loginUser models.LoginUser

	if c.ShouldBindJSON(&loginUser) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "error binding details to LoginUser",
		})
		return
	}

	// hash the password
	hpassword, herr := hashing.Hash(loginUser.Password)

	if herr != nil {
		panic("Error hashing the password")
	}

	newUser := models.User{
		Username: loginUser.Username,
		Password: hpassword,
		Admin:    false,
		OrgId:    orgid,
	}

	result := database.Db.Create(&newUser)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error creating new user",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"access_token": jwtAccessToken,
		"message":      "new user succesfully added",
	})

}

/* [ ] Function to delete a user from the organization */
func DeleteUser(c *gin.Context) {

	// extract the access token
	jwtAccessToken, oktoken := c.Get("access_token")

	if !oktoken {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting jwt access token from gin.context"})
		return
	}

	var deleteUser models.DeleteUser

	if c.ShouldBindJSON(&deleteUser) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "error binding details to deleteUser",
		})
		return
	}

	username := deleteUser.Username
	var findUser models.User

	// find the user with requested username
	result := database.Db.Where("username=?", username).Find(&findUser)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "requested user doesn't exist in database",
		})
		return
	}

	orgId, ok := c.Get("orgid")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error extracting orgid"})
		return
	}
	orgid := uint(orgId.(float64))

	// check if findUser's orgId is same as admin user
	if orgid != findUser.OrgId {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "requested user doesn't come under admin user's organization"})
		return
	}

	if database.Db.Delete(&models.User{Username: findUser.Username}).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting the requested user",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"access_token": jwtAccessToken,
		"message":      "requested user successfully deleted",
	})

}
