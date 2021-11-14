package rest

import (
	"dyn-https/blockchain/dynamic"
	"dyn-https/models"
	"dyn-https/util"
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

	jsonRPC := models.JSONRPC{}
	jsonRPC.ID, _ = util.RandomHashString(9)
	jsonRPC.Method = "audit"
	jsonRPC.JSONRPC = "1.0"
	jsonRPC.Params = make([]interface{}, 5)
	//audit add "hash_array" ( "account" ) ( "description" ) ( "algorithm" )
	jsonRPC.Params[0] = "add"
	jsonRPC.Params[1] = req.HashArray
	jsonRPC.Params[2] = req.Account
	jsonRPC.Params[3] = req.Description
	jsonRPC.Params[4] = req.Algorithm
	cmd, err := dynamic.GetNewRequest(jsonRPC)
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
	Id := c.Param("Id")
	jsonRPC := models.JSONRPC{}
	jsonRPC.ID, _ = util.RandomHashString(9)
	jsonRPC.Method = "audit"
	jsonRPC.JSONRPC = "1.0"
	jsonRPC.Params = make([]interface{}, 2)
	jsonRPC.Params[0] = "get"
	jsonRPC.Params[1] = Id
	cmd, err := dynamic.GetNewRequest(jsonRPC)
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

func (w *WebProxy) verifyaudit(c *gin.Context) {
	hash := c.Param("Hash")
	jsonRPC := models.JSONRPC{}
	jsonRPC.ID, _ = util.RandomHashString(9)
	jsonRPC.Method = "audit"
	jsonRPC.JSONRPC = "1.0"
	jsonRPC.Params = make([]interface{}, 2)
	jsonRPC.Params[0] = "verify"
	jsonRPC.Params[1] = hash
	cmd, err := dynamic.GetNewRequest(jsonRPC)
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
