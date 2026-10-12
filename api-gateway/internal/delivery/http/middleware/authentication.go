package middleware

import (
	"net/http"

	"github.com/LeonLow97/internal/contextstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type authenticationClaims struct {
	SessionID string `json:"session_id"`
	UserID    int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func Authentication(cookieName, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(cookieName)
		if err != nil {
			unauthorized(c)
			return
		}

		claims := &authenticationClaims{}
		token, err := jwt.ParseWithClaims(
			cookie.Value,
			claims,
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)
		if err != nil || token == nil || !token.Valid ||
			claims.UserID <= 0 || claims.SessionID == "" || claims.ExpiresAt == nil ||
			claims.Issuer != "authify-service" {
			unauthorized(c)
			return
		}

		ctx := contextstore.WithUserID(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
