package rest

import (
	"dyn-https/blockchain/dynamic"
	"dyn-https/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *WebProxy) addaudit(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	req := models.AuditAddRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	//audit add "hash_array" ( "account" ) ( "description" ) ( "algorithm" )
	cmd, _ := dynamic.NewRequest(`dynamic-cli audit add "` + req.HashArray + `" "` + req.Account +
		`"` + `" "` + req.Description + `"` + `" "` + req.Algorithm + `"`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (w *WebProxy) getaudit(c *gin.Context) {
	Id := c.Param("UserID")
	cmd, _ := dynamic.NewRequest(`dynamic-cli audit get "` + Id + `"`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	var result interface{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (w *WebProxy) verifyaudit(c *gin.Context) {
	hash := c.Param("Hash")
	cmd, _ := dynamic.NewRequest(`dynamic-cli audit verify "` + hash + `"`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	var result interface{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, result)
}
