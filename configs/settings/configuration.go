package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"dyn-https/models"
	"dyn-https/util"
)

// TODO: Support different ICE service authentication mechanisms

const (
	// ConfigurationFileName is the name of the configuration file
	ConfigurationFileName string = ".dynhttps.settings.json"
)

var homeDir string = ""
var pathSeperator string = ""

// HomeDir returns the dyn-https home directory
func HomeDir() string {
	return homeDir
}

// PathSeperator returns OS path seperator
func PathSeperator() string {
	return pathSeperator
}

// Configuration contains the main file settings used when the application launches
type Configuration struct {
	mut        *sync.RWMutex
	configFile models.ConfigurationFile
}

func isErr(e error) bool {
	if e != nil {
		return true
	}
	return false
}

func (c *Configuration) updateFile() {
	fileBytes, err := json.Marshal(&c.configFile)
	if err != nil {
		util.Error.Println("Configuration.updateFile() error after marshal configuration file data:", err)
	}
	fileName := (homeDir + ConfigurationFileName)
	err = ioutil.WriteFile(fileName, fileBytes, 0644)
	if err != nil {
		util.Error.Println("Configuration.updateFile() error after WriteFile: ", err)
	}
}

// ToJSON convert the Configuration struct to JSON
func (c *Configuration) ToJSON() models.ConfigurationFile {
	c.mut.Lock()
	defer c.mut.Unlock()
	// remove admin
	newAdmins := []models.Admin{}
	if c.configFile.Admins != nil {
		for _, admin := range c.configFile.Admins {
			admin.ExpectedHash = ""
			newAdmins = append(newAdmins, admin)
		}
	}
	c.configFile.Admins = newAdmins
	return c.configFile
}

func (c *Configuration) createDefault() {
	defaultWeb := models.DefaultWebServerConfig()
	defaultWalletStatus := models.DefaultWalletSetupStatus()
	c.configFile.WebServer = defaultWeb
	c.configFile.WalletStatus = defaultWalletStatus
	file, _ := json.Marshal(&c)
	err := ioutil.WriteFile(homeDir+ConfigurationFileName, file, 0644)
	if isErr(err) {
		util.Error.Println("Error writting default configuration file.")
	}
}

// Load reads the configuration file or loads default values
func (c *Configuration) Load(dir, seperator string) error {
	c.mut = new(sync.RWMutex)
	homeDir = dir
	pathSeperator = seperator
	_, errOpen := os.Open(homeDir + ConfigurationFileName)
	if isErr(errOpen) {
		util.Error.Println("Configuration file doesn't exist. Creating new configuration with default values.")
		c.createDefault()
	} else {
		file, errRead := ioutil.ReadFile(homeDir + ConfigurationFileName)
		if isErr(errRead) {
			c.createDefault()
			return nil
		}
		errUnmarshal := json.Unmarshal([]byte(file), &c.configFile)
		if isErr(errUnmarshal) {
			util.Error.Println("Error unmarshal configuration file. Overwritting file with default values.")
			c.createDefault()
		}
		if c.configFile.WebServer.AllowCIDR == "" && c.configFile.WebServer.BindAddress == "" && c.configFile.WebServer.ListenPort < 1 {
			c.configFile.WebServer = models.DefaultWebServerConfig()
			c.updateFile()
		} else {
			hasStatus := strings.Index(string(file), `"WalletStatus"`)
			if hasStatus < 0 {
				c.configFile.WalletStatus = models.DefaultWalletSetupStatus()
				c.updateFile()
			}
			if !util.IsValidCIDRList(c.configFile.WebServer.AllowCIDR) {
				return fmt.Errorf("Invalid Web Server allow CIDR: %v", c.configFile.WebServer.AllowCIDR)
			}
			if !util.IsValidIPAddress(c.configFile.WebServer.BindAddress) {
				return fmt.Errorf("Invalid Web Server bind IP address: %v", c.configFile.WebServer.BindAddress)
			}
			if c.configFile.WebServer.ListenPort < 1 {
				return fmt.Errorf("Invalid Web Server port: %v", c.configFile.WebServer.ListenPort)
			}
		}
		util.Info.Println("Configuration loaded successfully.")
	}
	if c.configFile.Admins == nil {
		// admin with password test
		//admin := models.Admin{UserName: "admin", ExpectedHash: "u9cYLNDulUiPGh5vP+DY+U7Q0U5NsdznE/6CoyMcUj0="}
		admins := []models.Admin{}
		c.configFile.Admins = admins
		c.updateFile()
	}
	return nil
}

func (c *Configuration) Admins() []models.Admin {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.Admins()
}

func (c *Configuration) AdminArrayToMap() map[string]string {
	admins := make(map[string]string)
	if c.configFile.Admins != nil && len(c.configFile.Admins) > 0 {
		for _, admin := range c.configFile.Admins {
			admins[admin.UserName] = admin.ExpectedHash
		}
	}
	return admins
}

func (c *Configuration) UpdateAdmins(admin models.Admin) {
	c.mut.Lock()
	defer c.mut.Unlock()
	index := -1
	for i, value := range c.configFile.Admins {
		if value.UserName == admin.UserName {
			index = i
			break
		}
	}
	if index != -1 {
		c.configFile.Admins[index] = admin
	} else {
		c.configFile.Admins = append(c.configFile.Admins, admin)
	}
	// update config file
	c.updateFile()
	return
}
