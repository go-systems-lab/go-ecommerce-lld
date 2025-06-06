package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key to use when setting gin context
const GinContextKey = "GinContextKey"

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
