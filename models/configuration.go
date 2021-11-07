package models

// ConfigurationFile stores the content of the dyn-https configuration file
// swagger:parameters models.ConfigurationFile
type ConfigurationFile struct {
	IceServers   []IceServerConfig `json:"IceServers"`
	WebServer    WebServerConfig
	WalletStatus WalletSetupStatus
}
