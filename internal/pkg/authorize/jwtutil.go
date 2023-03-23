package authorize

import (
	"example/auth/internal/pkg/models"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func init() {
	var loaderr = godotenv.Load()

	if loaderr != nil {
		log.Fatal("Error loading .env file")
	}
}

// custom claims for jwt
type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"admin"`
	OrgId    uint   `json:"orgid"`
	jwt.StandardClaims
}

func GetAccessTokenClaims(user models.User, expire int64) MyClaims {
	var accessTokenClaims = MyClaims{
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  user.Isadmin,
		OrgId:    user.OrgId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "melsonic",
			Subject:   "generate access token",
		},
	}

	return accessTokenClaims
}

func GetRefreshTokenClaims(user models.User, expire int64) MyClaims {
	var refreshTokenClaims = MyClaims{
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  user.Isadmin,
		OrgId:    user.OrgId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			Subject:   "generate refresh token",
		},
	}

	return refreshTokenClaims
}

func ExtractJwtClaims(jwtToken string) (map[string]interface{}, bool) {
	var claims map[string]interface{}
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySignature, nil
	})

	if err != nil {
		return claims, false
	}

	claims = token.Claims.(jwt.MapClaims)

	return claims, token.Valid
}

/* load the signature from .env file */
var mySignature = []byte(os.Getenv("JWT_SIGNATURE"))

func GenerateAccessToken(user models.User) (string, error) {
	var expire int64 = time.Now().Add(time.Hour).Unix()
	accessTokenClaims := GetAccessTokenClaims(user, expire)
	// create new JWT Token object with custom claims
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	// sign the token using a secret key
	SignedAccessToken, saerr := accessToken.SignedString(mySignature)

	return SignedAccessToken, saerr
}

/* function to generate refresh token */
func GenerateRefreshToken(user models.User) (string, int64, error) {
	var expire int64 = time.Now().Add(time.Hour * 24).Unix()
	refreshTokenClaims := GetRefreshTokenClaims(user, expire)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	SignedRefreshToken, srerr := refreshToken.SignedString(mySignature)
	return SignedRefreshToken, expire, srerr
}
