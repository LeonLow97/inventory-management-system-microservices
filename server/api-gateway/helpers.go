package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) retrieveUserIDFromToken(c *gin.Context) (int, error) {
	if userID, found := c.Get("userID"); found {
		userIDContext, err := strconv.Atoi(userID.(string))
		if err != nil {
			log.Println("Failed to convert userID to int in request context:", err)
			return 0, err
		}
		return userIDContext, nil
	} else {
		log.Println("UserID not found in jwt token claims")
		return 0, ErrMissingUserIDInJWTToken
	}
}
