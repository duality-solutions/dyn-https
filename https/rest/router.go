package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"dyn-https/blockchain/dynamic"
	"dyn-https/configs/settings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// WebBridgeRunner is used to run the node application.
type WebBridgeRunner struct {
	dynamicd      *dynamic.Dynamicd
	router        *gin.Engine
	configuration *settings.Configuration
	shutdownApp   *AppShutdown
	server        *http.Server
	mode          string
}

var runner WebBridgeRunner

// TODO: Add rate limitor
// TODO: Add custom logging
// TODO: Add authentication

// StartWebServiceRouter is used to setup the Rest server routes
func StartWebServiceRouter(c *settings.Configuration, d *dynamic.Dynamicd, a *AppShutdown, m string) {
	runner.configuration = c
	runner.dynamicd = d
	runner.shutdownApp = a
	runner.mode = m
	setupStatus, _, err := runner.GetWalletSetupInfo()
	if err == nil {
		runner.configuration.UpdateWalletSetup(*setupStatus)
	}
	startWebServiceRoutes()
}

func startWebServiceRoutes() {
	gin.SetMode(runner.mode)
	runner.router = gin.Default()
	runner.router.Use(AllowCIDR(runner.configuration.WebServer().AllowCIDR))
	setupAdminWebConsole()
	api := runner.router.Group("/api")
	version := api.Group("/v1")
	version.POST("/shutdown", runner.shutdown)
	version.GET("/overview", runner.overview)
	setupBlockchainRoutes(version)
	setupWalletRoutes(version)
	setupConfigRoutes(version)
	startGinGonic()
}

func startGinGonic() {
	runner.server = &http.Server{
		Addr:    runner.configuration.WebServer().AddressPortString(),
		Handler: runner.router,
	}
	// Start HTTP service
	if err := runner.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("ListenAndServe failed: %v", err))
	}
}

// RestartWebServiceRouter running service is shutdown and a new service is ran with a new configuration
func RestartWebServiceRouter() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := runner.server.Shutdown(ctx); err != nil {
		panic(fmt.Errorf("Server Shutdown: %v", err))
	}
	go startWebServiceRoutes()
}

func setupAdminWebConsole() {
	// Setup admin console web application
	runner.router.Use(static.Serve("/", static.LocalFile("./web/build", true)))
	runner.router.Use(static.Serve("/admin", static.LocalFile("./web/build", true)))
}

// @title WebBridge Restful API Documentation
// @version 1.0
// @description WebBridge Rest API discovery website.
// @termsOfService http://www.duality.solutions/webbridge/terms

// @contact.name API Support
// @contact.url http://www.duality.solutions/support
// @contact.email support@duality.solutions

// @license.name Duality
// @license.url https://github.com/duality-solutions/dyn-https/blob/master/LICENSE.md

// @host http://docs.dyn-https.duality.solutions
// @BasePath /api/v1
//func setupSwagger() {
//	address := runner.configuration.WebServer().AddressPortRawString() + "/swagger/doc.json"
//	url := ginSwagger.URL(address)
//	runner.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
//}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupBlockchainRoutes(currentVersion *gin.RouterGroup) {
	blockchain := currentVersion.Group("/blockchain")
	blockchain.POST("/jsonrpc", runner.handleJSONRPC)
	blockchain.GET("/", runner.getinfo)
	blockchain.GET("/sync", runner.syncstatus)
	blockchain.GET("/users", runner.users)
	blockchain.GET("/users/:UserID", runner.user)
	blockchain.GET("/groups", runner.groups)
	blockchain.GET("/groups/:GroupID", runner.group)
}

// TODO: follow https://rest.bitcoin.com for rest endpoints
func setupWalletRoutes(currentVersion *gin.RouterGroup) {
	wallet := currentVersion.Group("/wallet")
	wallet.GET("/", runner.walletinfo)
	wallet.GET("/mnemonic", runner.getmnemonic)
	wallet.POST("/mnemonic", runner.postmnemonic)
	wallet.PATCH("/unlock", runner.unlockwallet)
	wallet.PATCH("/lock", runner.lockwallet)
	wallet.PATCH("/encrypt", runner.encryptwallet)
	wallet.PATCH("/changepassphrase", runner.changepassphrase)
	wallet.GET("/users", runner.walletusers)
	wallet.GET("/groups", runner.walletgroups)
	wallet.GET("/links", runner.links)
	wallet.POST("/links/request", runner.linkrequest)
	wallet.POST("/links/accept", runner.linkaccept)
	wallet.POST("/links/message", runner.sendlinkmessage)
	wallet.GET("/links/message", runner.getlinkmessages)
	wallet.GET("/defaultaddress", runner.defaultaddress)
	wallet.GET("/transactions", runner.gettransactions)
	wallet.GET("/setup", runner.walletsetup)
	wallet.POST("/setup/backup", runner.setupmnemonicbackup)
}

func setupConfigRoutes(currentVersion *gin.RouterGroup) {
	config := currentVersion.Group("/config")
	config.GET("/", runner.config)
	config.GET("/web", runner.getwebserver)
	config.POST("/web", runner.updatewebserver)
	config.PUT("/web/restart", runner.restartwebserver)
}
