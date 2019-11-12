package beater

import (
	"os"
	"path/filepath"
)

// Returns Chromium history DB paths for each user on the machine
func getOperaPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Roaming", "Opera Software", "Opera Stable", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "com.operasoftware.Opera", "History")
		} else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "opera", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "opera.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}
