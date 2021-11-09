package rest

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuth(admins map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if ok {
			expectedHash := admins[username]
			if len(expectedHash) > 0 {
				hash := sha256.Sum256([]byte(username + password))
				exHash, err := base64.StdEncoding.DecodeString(expectedHash)
				if err == nil {
					hashMatch := (subtle.ConstantTimeCompare(hash[:], exHash[:]) == 1)
					if hashMatch {
						return
					}
				}
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
