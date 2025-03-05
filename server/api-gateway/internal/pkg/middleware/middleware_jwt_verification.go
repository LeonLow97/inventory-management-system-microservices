package middleware

import (
	"log"
	"strconv"
	"strings"

	"github.com/LeonLow97/internal/pkg/apierror"
	"github.com/LeonLow97/internal/pkg/contextstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
)

var skipPaths = map[string]struct{}{
	"/healthcheck": {},
	"/login":       {},
	"/signup":      {},
	"/logout":      {},
}

// JWTAuthMiddleware ensures that incoming requests have a valid JWT Token
// It extracts the token from the Authorization header, verifies its format,
// and validates the token's authenticity
func (m *Middleware) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip request path for this middleware
		if _, skipRequest := skipPaths[c.Request.URL.Path]; skipRequest {
			c.Next()
			return
		}

		// Retrieve the Authorization header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("[AUTH ERROR] Missing Authorization header in request")
			apierror.ErrUnauthorized.APIError(c, nil)
			return
		}

		// Ensure the Authorization header follows the "Bearer <token>" format
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			log.Println("[AUTH ERROR] Invalid Authorization header format. Expected 'Bearer <token>'")
			apierror.ErrUnauthorized.APIError(c, nil)
			return
		}

		if headerParts[0] != "Bearer" {
			log.Println("[AUTH ERROR] Missing 'Bearer' prefix in Authorization header")
			apierror.ErrUnauthorized.APIError(c, nil)
			return
		}

		// Parse and validate the JWT token
		jwtToken := headerParts[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("[JWT ERROR] Unexpected signing method")
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.cfg.AuthJWTToken.Secret), nil
		})

		// Handle token validation failure
		if err != nil || !token.Valid {
			log.Printf("[AUTH ERROR] Invalid or expired token: %v\n", err)
			apierror.ErrUnauthorized.APIError(c, nil)
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[JWT ERROR] Failed to parse token claims")
			apierror.ErrInternalServerError.APIError(c, nil)
			return
		}

		// Retrieve the issuer (user ID) from the token claims
		issuer, ok := claims["iss"]
		if !ok {
			log.Println("[JWT ERROR] Missing 'iss' (issuer) field in token claims")
			apierror.ErrInternalServerError.APIError(c, nil)
			return
		}

		// Convert issuer to string (if not already) and then to int
		userIDStr, ok := issuer.(string)
		if !ok {
			log.Println("[JWT ERROR] 'iss' (issuer) is not a string")
			apierror.ErrInternalServerError.APIError(c, nil)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println("[JWT ERROR] Failed to convert 'iss' to integer:", err)
			apierror.ErrInternalServerError.APIError(c, nil)
			return
		}

		if userID <= 0 {
			log.Printf("User ID %d is invalid", userID)
			apierror.ErrUnauthorized.APIError(c, nil)
			return
		}

		// gRPC metadata for outgoing gRPC calls
		md := metadata.New(map[string]string{
			"user_id": userIDStr,
		})

		// Store the metadata in Gin context
		contextstore.InjectGRPCMetadataIntoContext(c, md)

		c.Next()
	}
}
