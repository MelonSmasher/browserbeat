package beater

import (
	"os"
	"path/filepath"
)

// Returns Chrome history DB paths for each user on the machine
func getChromePaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Google", "Chrome", "Default", "History")
		} else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "google-chrome", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "chrome.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Chrome Beta history DB paths for each user on the machine
func getChromeBetaPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "google-chrome-beta", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "chrome-beta.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Chrome Unstable history DB paths for each user on the machine
func getChromeDevPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "google-chrome-unstable", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "chrome-unstable.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Chrome Canary history DB paths for each user on the machine
func getChromeCanaryPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Google", "Chrome SxS", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Google", "Chrome Canary", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "chrome-canary.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}
