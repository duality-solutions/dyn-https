package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dyn-https/blockchain/dynamic"
	"dyn-https/models"

	"github.com/gin-gonic/gin"
)

func (w *WebProxy) users(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli getusers`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	var result interface{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (w *WebProxy) user(c *gin.Context) {
	userID := c.Param("UserID")
	cmd := `dynamic-cli getuserinfo "` + userID + `"`
	strCommand, err := dynamic.NewRequest(cmd)
	if err != nil {
		strErrMsg := fmt.Sprintf("NewRequest cmd(%v), error: %v", cmd, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func (w *WebProxy) walletusers(c *gin.Context) {
	strCommand, _ := dynamic.NewRequest(`dynamic-cli mybdapaccounts`)
	response, _ := <-w.dynamicd.ExecCmdRequest(strCommand)
	myAccounts := make(map[string]models.Account)
	err := json.Unmarshal([]byte(response), &myAccounts)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	myUsers := make(map[string]models.Account)
	for i, account := range myAccounts {
		if account.ObjectType == "User Entry" {
			myUsers[i] = account
		}
	}
	c.JSON(http.StatusOK, myUsers)
}
