package rest

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"dyn-https/blockchain/dynamic"
	"dyn-https/configs/settings"
	_ "dyn-https/docs" // used for Swagger documentation

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// WebProxy is used to run the node application.
type WebProxy struct {
	dynamicd      *dynamic.Dynamicd
	router        *gin.Engine
	configuration *settings.Configuration
	shutdownApp   *AppShutdown
	server        *http.Server
	mode          string
	Admins        map[string]string
}

var runner WebProxy

// TODO: Add rate limitor
// TODO: Add custom logging
// TODO: Add admin username and password commands

// StartWebServiceRouter is used to setup the Rest server routes
func StartWebServiceRouter(c *settings.Configuration, d *dynamic.Dynamicd, a *AppShutdown, m string) {
	runner.configuration = c
	runner.dynamicd = d
	runner.shutdownApp = a
	runner.mode = m
	runner.Admins = c.AdminArrayToMap()
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
	runner.router.Use(BasicAuth(runner.Admins))
	setupAdminWebConsole()
	api := runner.router.Group("/api")
	version := api.Group("/v1")
	version.POST("/shutdown", runner.shutdown)
	version.GET("/overview", runner.overview)
	setupBlockchainRoutes(version)
	setupWalletRoutes(version)
	setupConfigRoutes(version)
	setupSwagger()
	startGinGonic()
}

func setupLetsEncrypt(r *gin.Engine, domains []string) {
	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	log.Fatal(autotls.Run(r, domains...))
}

func getSelfSignedCert() tls.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		// todo: get the cert values from config
		Subject: pkix.Name{
			Organization:  []string{"Duality Blockchain Solutions"},
			Country:       []string{"USA"},
			Province:      []string{"Texas"},
			Locality:      []string{"San Antonio"},
			StreetAddress: []string{"110 E Houston St, 7th Floor"},
			PostalCode:    []string{"78205"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	// Sign the certificate
	certificate, _ := x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)

	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// Generate a key pair from your pem-encoded cert and key ([]byte).
	x509Cert, _ := tls.X509KeyPair(certBytes, keyBytes)
	return x509Cert
}

func startGinGonic() {
	switch strings.ToLower(runner.mode) {
	case "debug":
		cert := getSelfSignedCert()
		runner.server = &http.Server{
			Addr:    runner.configuration.WebServer().AddressPortString(),
			Handler: runner.router,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
		// Start HTTPS service
		if err := runner.server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("ListenAndServe failed: %v", err))
		}
	case "release":
		// todo: get the domains from config
		domains := []string{"dynhttps.dualityblocks.com", "dyn.pix.dualityblocks.com"}
		setupLetsEncrypt(runner.router, domains)
		runner.server = &http.Server{
			Addr:    runner.configuration.WebServer().AddressPortString(),
			Handler: runner.router,
		}
		// Start HTTPS service
		if err := runner.server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("ListenAndServe failed: %v", err))
		}
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

// @title DYN HTTPS Restful API Documentation
// @version 1.0
// @description DYN HTTPS Rest API discovery website.
// @termsOfService http://www.duality.solutions/dynhttps/terms

// @contact.name API Support
// @contact.url http://www.duality.solutions/support
// @contact.email support@duality.solutions

// @license.name Duality
// @license.url https://github.com/duality-solutions/dyn-https/blob/master/LICENSE.md

// @host http://docs.dyn-https.duality.solutions
// @BasePath /api/v1

func setupSwagger() {
	address := runner.configuration.WebServer().AddressPortRawString() + "/swagger/doc.json"
	url := ginSwagger.URL(address)
	runner.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

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
	blockchain.POST("/audit", runner.addaudit)
	blockchain.GET("/audit/:Id", runner.getaudit)
	blockchain.GET("/audit/verify/:Hash", runner.verifyaudit)
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
