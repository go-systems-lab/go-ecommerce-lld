package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
)

func AuthorizeJWT(jwtService account.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("token")
		if err != nil || authCookie == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := jwtService.ValidateToken(authCookie)
		if err != nil {
			// Token is invalid => treat as anonymous or invalid user
			c.Set("userID", "")
			c.Next()
			return
		}

		// Token is valid => set user info
		if claims, ok := token.Claims.(*account.JWTCustomClaims); ok && token.Valid {
			c.Set("userID", claims.UserID)
		} else {
			c.Set("userID", "")
		}

		// Continue the request
		c.Next()
	}
}
