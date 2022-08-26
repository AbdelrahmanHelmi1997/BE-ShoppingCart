package middleware

import (
	"fmt"
	"net/http"

	"AmzonElGalaba/Helper"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

func (s *Auth) Authentication(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
		c.Abort()
		return
	}

	claims, err := Helper.ValidateToken(clientToken)
	if claims.Role_Type != "admin" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("don't Have access")})
		c.Abort()
		return
	}
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

}
