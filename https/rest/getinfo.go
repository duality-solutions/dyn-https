package rest

import (
	"encoding/json"
	"net/http"

	"dyn-https/blockchain/dynamic"
	"dyn-https/models"

	"github.com/gin-gonic/gin"
)

func (w *WebProxy) getinfo(c *gin.Context) {
	var info models.GetInfoData
	req, _ := dynamic.NewRequest("dynamic-cli getinfo")
	err := json.Unmarshal([]byte(<-w.dynamicd.ExecCmdRequest(req)), &info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (w *WebProxy) syncstatus(c *gin.Context) {
	var status models.SyncStatus
	req, _ := dynamic.NewRequest("dynamic-cli syncstatus")
	err := json.Unmarshal([]byte(<-w.dynamicd.ExecCmdRequest(req)), &status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, status)
}
