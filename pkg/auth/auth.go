package auth

import (
	"aspire-assignment/pkg/config"
	e "aspire-assignment/pkg/errors"
	"log"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	jwtKey []byte
)

func InitAuth() {
	confi := config.GetConfig()
	jwtKey = []byte(confi.GetString("auth.key"))
}

func GenerateJWT(payload Token) (string, error) {
	expirationTime := payload.Exp
	claims := &Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response AuthResponse
		)
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" || !strings.Contains(tokenString, "Bearer") {
			log.Println("auth token incorrect")
			response.Status = false
			response.Message = "invalid jwt token format"
			response.Errors = append(response.Errors, e.ErrorInfo[e.UnAuthorized].GetErrorDetails("invalid auth token"))
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if claims.Payload.UserId == 0 {
			log.Println("claims are unavailable from the token")
			response.Status = false
			response.Message = "Authentication failed"
			response.Errors = append(response.Errors, e.ErrorInfo[e.UnAuthorized].GetErrorDetails("invalid auth token"))
			c.JSON(http.StatusForbidden, response)
			c.Abort()
			return
		}

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				response.Status = false
				response.Message = "invalid auth token signature"
				response.Errors = append(response.Errors, e.ErrorInfo[e.UnAuthorized].GetErrorDetails("invalid auth token signature"))
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
			response.Status = false
			response.Message = "invalid jwt token"
			response.Errors = append(response.Errors, e.ErrorInfo[e.UnAuthorized].GetErrorDetails("invalid jwt token"))
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		if !token.Valid {
			response.Status = false
			response.Message = "expired jwt token"
			response.Errors = append(response.Errors, e.ErrorInfo[e.UnAuthorized].GetErrorDetails("expired token"))
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set(config.USERID, claims.Payload.UserId)
		c.Set(config.USERNAME, claims.Payload.UserName)
		c.Set(config.USERTYPE, claims.Payload.UserType)

		//block all non admin/ path calls for ADMIN and admin/ path calls for CUSTOMER
		if (claims.Payload.UserType == config.ADMIN && !strings.Contains(c.FullPath(), "/admin/")) ||
			(claims.Payload.UserType == config.CUSTOMER && strings.Contains(c.FullPath(), "/admin/")) {
			response.Status = false
			response.Message = "access not allowed"
			response.Errors = append(response.Errors, e.ErrorInfo[e.UnAuthorized].GetErrorDetails("access not allowed"))
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		c.Next()
	}
}
