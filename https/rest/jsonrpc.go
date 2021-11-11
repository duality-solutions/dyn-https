package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dyn-https/blockchain/dynamic"
	"dyn-https/models"

	"github.com/gin-gonic/gin"
)

func (w *WebProxy) handleJSONRPC(c *gin.Context) {
	reqInput := models.JSONRPC{}
	err := json.NewDecoder(c.Request.Body).Decode(&reqInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	req, err := dynamic.GetNewRequest(reqInput)
	if err != nil {
		strErrMsg := fmt.Sprintf("GetNewRequest error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	response, _ := <-w.dynamicd.ExecCmdRequest(req)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Result unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, result)
}
