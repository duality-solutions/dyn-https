package configs

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"dyn-https/configs/settings"
	"dyn-https/models"
)

func AdminUserCommand(c *settings.Configuration, cmd string) (string, error) {
	cmd = strings.Replace(cmd, "useradd ", "", -1)
	spaceindex := strings.Index(cmd, " ")
	if spaceindex <= 0 {
		return "", errors.New("Invalid paramaters")
	}
	params := strings.Split(cmd, " ")
	if len(params) != 2 {
		return "", fmt.Errorf("Invalid number of paramaters %d", len(params))
	}
	username := params[0]
	password := params[1]
	if len(password) < 8 {
		return username, fmt.Errorf("Invalid password length %d. Minimum password must be at least 8 charactors.", len(password))
	}
	// todo: check password complexity
	hash := sha256.Sum256([]byte(username + password))
	expectedHash := base64.StdEncoding.EncodeToString(hash[:])
	admin := models.Admin{UserName: username, ExpectedHash: expectedHash}
	if len(admin.ExpectedHash) == 0 || len(admin.UserName) == 0 {
		return username, fmt.Errorf("Invalid number of paramaters %d", len(params))
	}
	c.UpdateAdmins(admin)
	return username, nil
}
