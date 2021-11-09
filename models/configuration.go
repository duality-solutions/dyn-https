package models

// ConfigurationFile stores the content of the dyn-https configuration file
// swagger:parameters models.ConfigurationFile
type ConfigurationFile struct {
	Admins       []Admin
	WebServer    WebServerConfig
	WalletStatus WalletSetupStatus
}
