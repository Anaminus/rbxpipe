package main

import (
	"os"
	path "path/filepath"
)

// A list of valid studio executable names.
var Executables = []string{
	"RobloxStudioBeta.exe",
}

// Location of plugin settings.
var PluginSettings = path.Join(localAppData(), `Roblox\InstalledPlugins\0\settings.json`)

// Location of Roblox builds for all users.
var AllUsers = path.Join(programFiles(), `Roblox\Versions`)

// Location of Roblox builds for the current user.
var CurrentUser = path.Join(localAppData(), `Roblox\Versions`)

func localAppData() string {
	lappdata := os.Getenv("LOCALAPPDATA")
	if _, err := os.Stat(lappdata); lappdata == "" || err != nil {
		userProfile := os.Getenv("USERPROFILE")
		lappdata = path.Join(userProfile, `AppData\Local`)
		if _, err := os.Stat(lappdata); lappdata == "" || err != nil {
			lappdata = path.Join(userProfile, `Local Settings\Application Data`)
		}
	}
	return lappdata
}

func programFiles() string {
	programFiles := `C:\Program Files (x86)`
	if _, err := os.Stat(programFiles); err != nil {
		programFiles = `C:\Program Files`
	}
	return programFiles
}
