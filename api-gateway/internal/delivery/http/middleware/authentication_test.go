package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LeonLow97/internal/contextstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	testCookieName = "auth-token"
	testJWTSecret  = "test-secret"
)

func signedTestToken(t *testing.T, method jwt.SigningMethod, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}
	return signedToken
}

func validTestClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"user_id":    42,
		"session_id": "session-123",
		"iss":        "authify-service",
		"exp":        time.Now().Add(time.Hour).Unix(),
	}
}

func TestAuthentication(t *testing.T) {
	validToken := signedTestToken(t, jwt.SigningMethodHS256, validTestClaims(), testJWTSecret)
	wrongSignatureToken := signedTestToken(t, jwt.SigningMethodHS256, validTestClaims(), "another-secret")
	wrongIssuerClaims := validTestClaims()
	wrongIssuerClaims["iss"] = "another-service"
	expiredClaims := validTestClaims()
	expiredClaims["exp"] = time.Now().Add(-time.Hour).Unix()
	missingUserIDClaims := validTestClaims()
	delete(missingUserIDClaims, "user_id")
	zeroUserIDClaims := validTestClaims()
	zeroUserIDClaims["user_id"] = 0
	missingSessionClaims := validTestClaims()
	delete(missingSessionClaims, "session_id")
	missingExpiryClaims := validTestClaims()
	delete(missingExpiryClaims, "exp")

	testCases := []struct {
		name       string
		token      string
		wantStatus int
		wantUserID int64
	}{
		{name: "valid token", token: validToken, wantStatus: http.StatusNoContent, wantUserID: 42},
		{name: "missing cookie", wantStatus: http.StatusUnauthorized},
		{name: "invalid token", token: "not-a-jwt", wantStatus: http.StatusUnauthorized},
		{name: "wrong signature", token: wrongSignatureToken, wantStatus: http.StatusUnauthorized},
		{name: "wrong signing algorithm", token: signedTestToken(t, jwt.SigningMethodHS512, validTestClaims(), testJWTSecret), wantStatus: http.StatusUnauthorized},
		{name: "wrong issuer", token: signedTestToken(t, jwt.SigningMethodHS256, wrongIssuerClaims, testJWTSecret), wantStatus: http.StatusUnauthorized},
		{name: "expired token", token: signedTestToken(t, jwt.SigningMethodHS256, expiredClaims, testJWTSecret), wantStatus: http.StatusUnauthorized},
		{name: "missing user id", token: signedTestToken(t, jwt.SigningMethodHS256, missingUserIDClaims, testJWTSecret), wantStatus: http.StatusUnauthorized},
		{name: "zero user id", token: signedTestToken(t, jwt.SigningMethodHS256, zeroUserIDClaims, testJWTSecret), wantStatus: http.StatusUnauthorized},
		{name: "missing session id", token: signedTestToken(t, jwt.SigningMethodHS256, missingSessionClaims, testJWTSecret), wantStatus: http.StatusUnauthorized},
		{name: "missing expiry", token: signedTestToken(t, jwt.SigningMethodHS256, missingExpiryClaims, testJWTSecret), wantStatus: http.StatusUnauthorized},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(Authentication(testCookieName, testJWTSecret))
			var gotUserID int64
			router.GET("/protected", func(c *gin.Context) {
				gotUserID, _ = contextstore.UserIDFromContext(c.Request.Context())
				c.Status(http.StatusNoContent)
			})

			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if testCase.token != "" {
				request.AddCookie(&http.Cookie{Name: testCookieName, Value: testCase.token})
			}
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, testCase.wantStatus)
			}
			if testCase.wantStatus == http.StatusNoContent && gotUserID != testCase.wantUserID {
				t.Errorf("context user ID = %d, want %d", gotUserID, testCase.wantUserID)
			}
		})
	}
}
