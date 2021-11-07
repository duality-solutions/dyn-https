package dynamic

import (
	"os/exec"
	"path/filepath"
	"strings"

	"dyn-https/util"
)

func downloadBinaries(_os, dynDir, dynamicName, cliName, archiveExt string) error {
	if !util.FileExists(dynDir+dynamicdName) || !util.FileExists(dynDir+cliName) {
		util.Info.Println(dynamicdName, "or", cliName, "not found. Downloading from Git repo.")
		binPath := dynDir + "dynamic." + archiveExt
		if !util.FileExists(binPath) {
			binaryURL := binaryRepo + "/" + binaryReleasePath + "/v" + binaryVersion + "/Dynamic-" + binaryVersion + "-" + _os + "-" + arch + "." + archiveExt
			util.Info.Println("Downloading binaries:", binaryURL)
			err := util.DownloadFile(binPath, binaryURL)
			if err != nil {
				util.Error.Println("Binary download failed.", err)
				return err
			}
		} else {
			util.Info.Println("Compressed binary found")
		}
		var binaries []string
		// move keep files to const and eventually a yaml file
		var keepFiles = []string{"bin/dynamicd", "bin/dynamic-cli"}
		var errDecompress error
		// todo test if Windows works with keep files
		if _os == "Windows" {
			// unzip archive file
			binaries, errDecompress = util.Unzip(binPath, dynDir)
			if errDecompress != nil {
				util.Error.Println("Error unzipping file.", binPath, errDecompress)
				return errDecompress
			}
		} else {
			// Extract tar.gz archive file
			binaries, errDecompress = util.ExtractTarGz(binPath, dynDir, keepFiles)
			if errDecompress != nil {
				util.Error.Println("Error decompressing file. binPath: ", binPath, " dynDir: ", dynDir, " Error: ", errDecompress)
				return errDecompress
			}
		}
		util.Info.Println("downloadBinaries keep file", keepFiles)
		for _, v := range binaries {
			util.Info.Println("downloadBinaries keep file", v)
			if _os != "Windows" {
				util.Info.Println("downloadBinaries chmod ", dynDir+filepath.Base(v), " GetFileNameFromPath ", filepath.Base(v))
				cmd := exec.Command("chmod", "+x", dynDir+filepath.Base(v))
				errRun := cmd.Run()
				if errRun != nil {
					util.Error.Println("Error running chmod for", dynDir+v, errRun)
				}
			}
		}
		// clean up extract directory
		dirs, errDirs := util.ListDirectories(dynDir)
		if errDirs != nil {
			util.Error.Println("Error listing directories", errDirs)
			return errDirs
		}
		for _, v := range dirs {
			if !strings.HasPrefix(v, ".dynamic") {
				util.Info.Println("Deleting directory", dynDir+v)
				errDeleteDir := util.DeleteDirectory(dynDir + v)
				if errDeleteDir != nil {
					util.Error.Println("Error deleting directory", v, errDeleteDir)
					return errDeleteDir
				}
			}
		}
		// clean up archive file
		util.Info.Println("Cleaning up... Removing unneeded files and directories.")
		if util.FileExists(binPath) {
			util.Info.Println("Deleting zip file", binPath)
			errDelete := util.DeleteFile(binPath)
			if errDelete != nil {
				util.Error.Println("Error deleting binary archive file", binPath, errDelete)
				return errDelete
			}
		}
	} else {
		util.Info.Println("Binaries found", dynamicdName, cliName)
	}
	return nil
}
