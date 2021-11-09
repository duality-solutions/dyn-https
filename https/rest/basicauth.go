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
				// add username salt to slow down brute force password hash cracking
				// ToDo: consider using a better password hash algorithm like Argon2i
				hash := sha256.Sum256([]byte(username + password))
				exHash, err := base64.StdEncoding.DecodeString(expectedHash)
				if err == nil {
					// To prevent timing attacks based on error length
					hashMatch := (subtle.ConstantTimeCompare(hash[:], exHash[:]) == 1)
					if hashMatch {
						return
					}
				}
			} else {
				// Hash, decode and compare to prevent timing attacks based on username not found
				hash := sha256.Sum256([]byte(username + password))
				exHash, _ := base64.StdEncoding.DecodeString("u9cYLNDulUiPGh5vP+DY+U7Q0U5NsdznE/6CoyMcUj0=")
				subtle.ConstantTimeCompare(hash[:], exHash[:])
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
