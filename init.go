package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"time"

	"dyn-https/blockchain/dynamic"
	"dyn-https/configs/settings"
	"dyn-https/https/rest"
	"dyn-https/util"
)

const (
	// DefaultName application name
	DefaultName string = "dyn-https"
)

var config settings.Configuration
var debug = false
var test = false
var walletpassphrase = ""
var testCreateOffer = false
var testWaitForOffer = false

// Init is used to begin all dyn-https tasks
func Init(version, githash string) error {
	running, pid, err := util.FindProcess(DefaultName)
	if running {
		if err == nil && pid > 0 {
			return fmt.Errorf("%v process (%v) found running. Can only run one instance at a time", DefaultName, pid)
		}
		return fmt.Errorf("%v process (%v) found running. Can only run one instance at a time. Error: %v", DefaultName, pid, err)
	}
	args := os.Args[1:]
	if len(args) > 0 {
		for _, v := range args {
			switch v {
			case "-debug":
				debug = true
			case "-test":
				test = true
			}
		}
	}
	if test || debug {
		/*
			if !testauth.TestCerts() {
				return fmt.Errorf("%v could not load test certificates", DefaultName)
			}
			if !testauth.TestExistingSignature() {
				return fmt.Errorf("%v could run existing key signature test", DefaultName)
			}
			if !testauth.TestNewKeySignature() {
				return fmt.Errorf("%v could run new key signature test", DefaultName)
			}
		*/
	}
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	pathSeperator := ""
	if runtime.GOOS == "windows" {
		pathSeperator = `\\`
		homeDir += pathSeperator + `.` + DefaultName + pathSeperator
	} else {
		pathSeperator = `/`
		homeDir += pathSeperator + `.` + DefaultName + pathSeperator
	}
	// initilize debug.log file
	util.InitDebugLogFile(debug, homeDir)
	util.Info.Println("Version:", version, "Hash", githash)
	util.Info.Println("OS: ", runtime.GOOS)
	if debug {
		util.Info.Println("Running", DefaultName, "in debug log mode.")
		util.Info.Println("Args", args)
	}
	if test {
		util.Info.Println("Running %v in test mode.", DefaultName)
	}
	err = config.Load(homeDir, pathSeperator)
	if err != nil {
		util.Info.Println("Error loading configuration file. %v", err)
	}
	util.Info.Println("Config", config)

	if !test {
		dynamicd, proc, err := dynamic.FindDynamicdProcess(false, time.Second*1)
		if err == nil {
			// kill existing dynamicd process
			util.Warning.Println("dynamicd daemon already running. Attempting to kill the process.")
			err = proc.Kill()
			if err != nil {
				return fmt.Errorf("Fatal error, dynamicd daemon process (%v) is running but can't be stopped %v", proc.Pid, err)
			}
			time.Sleep(time.Second * 5)
		}
		dynamicd, proc, err = dynamic.FindDynamicdProcess(true, time.Second*30)
		if proc != nil {
			util.Info.Println("Running dynamicd process found Pid", proc.Pid)
		} else {
			return fmt.Errorf("Fatal error starting dynamicd daemon %v ", err)
		}
		// make sure wallet is created
		dynamicd.WaitForWalletCreated()
		status, errStatus := dynamicd.GetSyncStatus()
		if errStatus != nil {
			return fmt.Errorf("GetSyncStatus %v", errStatus)
		}
		util.Info.Println("dynamicd running... Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")

		info, errInfo := dynamicd.GetInfo()
		if errInfo != nil {
			return fmt.Errorf("GetInfo %v", errInfo)
		}
		util.Info.Println("dynamic connections", info.Connections)

		acc, errAccounts := dynamicd.GetMyAccounts(time.Second * 120)
		if errAccounts != nil {
			util.Error.Println("GetActiveLinks error", errAccounts)
		} else {
			for i, account := range *acc {
				util.Info.Println("Account", i+1, account.CommonName, account.ObjectFullPath, account.WalletAddress, account.LinkAddress)
			}
		}
		var mode string = "release"
		if debug {
			mode = "debug"
		}
		// Create ShutdownApp stuct
		stopWatcher := make(chan struct{})
		dynamic.WatchProcess(stopWatcher, 10, walletpassphrase)
		var sync = false
		closeApp := make(chan struct{})
		shutdown := rest.AppShutdown{
			Close:       &closeApp,
			StopWatcher: &stopWatcher,
			Dynamicd:    dynamicd,
		}
		// Start Gin web services
		go rest.StartWebServiceRouter(&config, dynamicd, &shutdown, mode)

		al, errLinks := dynamicd.GetActiveLinks(time.Second * 120)
		if errLinks != nil {
			util.Error.Println("GetActiveLinks error", errLinks)
		} else {
			util.Info.Printf("Found %v links\n", len(al.Links))
		}

		appCommandLoop(acc, al, &shutdown, dynamicd, status, sync, &config)

		for {
			select {
			case <-*shutdown.Close:
				util.Info.Println("Shutdown close trigger.")
				return nil
			}
		}
	}
	return nil
}
