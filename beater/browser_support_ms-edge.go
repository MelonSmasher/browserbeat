package beater

import (
	"os"
	"path/filepath"
)

// Returns Edge history DB paths for each user on the machine
func getEdgePaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Microsoft Edge", "Default", "History")
		}
		/*else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "microsoft-edge", "Default", "History")
		}*/
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "edge.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Edge Beta history DB paths for each user on the machine
func getEdgeBetaPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Microsoft", "Edge Beta", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Microsoft Edge Beta", "Default", "History")
		}
		/*else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "microsoft-edge-beta", "Default", "History")
		}*/
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "edge-beta.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Edge Unstable history DB paths for each user on the machine
func getEdgeDevPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Microsoft", "Edge Dev", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Microsoft Edge Dev", "Default", "History")
		}
		/*else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "microsoft-edge-dev", "Default", "History")
		}*/
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "edge-dev.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Edge Canary history DB paths for each user on the machine
func getEdgeCanaryPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Microsoft", "Edge SxS", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Microsoft Edge Canary", "Default", "History")
		}
		/*else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "microsoft-edge-canary", "Default", "History")
		}*/
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "edge-canary.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}
