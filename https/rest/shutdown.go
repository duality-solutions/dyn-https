package rest

import (
	"net/http"
	"time"

	"dyn-https/blockchain/dynamic"
	"dyn-https/util"

	"github.com/gin-gonic/gin"
)

// AppShutdown stores the channels and commands needed to safely shutdown the DYN-HTTPS service
type AppShutdown struct {
	Close       *chan struct{}
	Dynamicd    *dynamic.Dynamicd
	StopWatcher *chan struct{}
}

// ShutdownAppliction safely shuts down process watcher, and the Dynamic daemon before exiting
func (a *AppShutdown) ShutdownAppliction() {
	close(*a.StopWatcher)

	// Stop dynamicd
	reqStop, _ := dynamic.NewRequest("dynamic-cli stop")
	respStop, _ := util.BeautifyJSON(<-a.Dynamicd.ExecCmdRequest(reqStop))
	util.Info.Println(respStop)
	time.Sleep(time.Second * 5)
	util.Info.Println("Looking for dynamicd process pid", a.Dynamicd.Cmd.Process.Pid)
	util.WaitForProcessShutdown(a.Dynamicd.Cmd.Process, time.Second*240)
	util.Info.Println("Exit.")
	util.EndDebugLogFile(30)
	close(*a.Close)
}

func (w *WebBridgeRunner) shutdown(c *gin.Context) {
	if w.shutdownApp != nil {
		w.shutdownApp.ShutdownAppliction()
		c.JSON(http.StatusOK, gin.H{"result": "WebBridge is shutting down."})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration variable shutdownApp is null."})
	}
}
