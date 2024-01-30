package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct{
	Name string
	jwt.StandardClaims
}

func GetToken(key string) gin.HandlerFunc{
	return func(c *gin.Context) {

		claims := &Claims{
		Name: "ODDS",
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(time.Minute * 5).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims,)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token" : signedToken,
	})
	}
	
}

func AuthMiddleware(key string) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")

		jwtToken , err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil || !jwtToken.Valid{
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims , ok := jwtToken.Claims.(*Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("Name", claims.Name,)
		ctx.Next()
		
	}
}