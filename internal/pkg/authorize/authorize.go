package authorize

import (
	"net/http"
	"time"

	"example/auth/internal/pkg/models"

	"example/auth/internal/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/* function to store Refresh Token into the database */
func StoreRefreshToken(db *gorm.DB, username string, refresh_token string, expire int64) bool {
	var RefreshToken models.JwtRefreshToken

	rows := db.Where("username=?", username).Find(&RefreshToken)

	// if already present in the refresh token update
	if rows.RowsAffected != 0 {
		result := db.Model(RefreshToken).UpdateColumn("refreshToken", refresh_token)
		return result.Error == nil
	}

	// if not present already create a new row
	RefreshToken = models.JwtRefreshToken{
		Username:     username,
		RefreshToken: refresh_token,
		Expire:       expire,
	}

	result := db.Create(&RefreshToken)

	// if already present then update
	return result.Error == nil
}

/* middleware function to extract jwt access token */
func ExtractJwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract the jwt token from request header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "user needs to be logged in...",
			})
			return
		}

		// remove the beared
		jwtToken := authHeader[7:]

		/* CHECK IF TOKEN IS BLACKLISTED */
		// declare a blacklisted token
		var blacklistedAccess models.JwtBlackListedToken
		result := database.Db.Where("accessToken=?", jwtToken).Find(&blacklistedAccess)

		if result.RowsAffected != 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "please log in, access token blacklisted"})
			return
		}

		claims, isTokenValid := ExtractJwtClaims(jwtToken)

		if len(claims) == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error extracting the token"})
			return
		}

		if !isTokenValid {
			jwtToken = GenerateAccessFromRefresh(claims["username"].(string))
			if jwtToken == "" {
				c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "error creating new token, may be refresh token expired"})
				return
			}
			claims, _ = ExtractJwtClaims(jwtToken) /* since validity will always be true for newly created token */
		}

		c.Set("access_token", jwtToken)
		c.Set("user", claims)
		c.Next()
	}
}

/* Function to generate new Access Token from Refresh */
func GenerateAccessFromRefresh(username string) string {
	var refresh_token models.JwtRefreshToken
	var user models.User

	// first check if the refresh token is valid or not
	result := database.Db.Where("username=?", username).Find(&refresh_token)

	if result.Error != nil || (time.Now().Unix() > refresh_token.Expire) {
		// delete the refresh token from the database
		database.Db.Where("username=?", username).Delete(&refresh_token)
		return ""
	}

	resultUser := database.Db.Where("username=?", username).Find(&user)

	if resultUser.Error != nil {
		return ""
	}

	access_token, _, _ := GenerateAccessToken(user)

	return access_token
}
