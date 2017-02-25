package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
	"time"
)

var JWT_SECRET []byte = []byte(os.Getenv("JWT_SECRET"))

type PublicUser struct {
	Id   bson.ObjectId
	Name string
	Role string
}

func MarshalToken(name, id, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"name": name,
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(JWT_SECRET)
}

func unMarshalToken(tokenStr string) (*PublicUser, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Problem with the JWT")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &PublicUser{
			bson.ObjectIdHex(claims["id"].(string)),
			claims["name"].(string),
			claims["role"].(string),
		}, nil
	} else {
		return nil, fmt.Errorf("Problem with the JWT")
	}
}

func IsAuthentified(c *gin.Context) {
	auth, err := c.Cookie("auth")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}

	pu, err := unMarshalToken(auth)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/users/login")
		c.Abort()
		return
	}

	c.Set("user", pu)
	c.Next()
}

func IsProf(status bool, responseType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*PublicUser)

		if user.Role == "teacher" && status || user.Role != "teacher" && !status {
			c.Next()
		} else {
			c.Abort()

			if responseType == "json" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
			} else {
				c.Redirect(http.StatusSeeOther, "/groups")
			}
		}
	}
}

func DeleteAuthCookie(c *gin.Context) {
	c.SetCookie("auth", "", -1, "", "", os.Getenv("env") == "prod", true)
	c.Next()
}
