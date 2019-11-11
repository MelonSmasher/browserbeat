package beater

import (
	"os"
	"path/filepath"
)

// Returns Chromium history DB paths for each user on the machine
func getVivaldiPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Vivaldi", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Vivaldi", "Default", "History")
		} else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "vivaldi-snapshot", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "vivaldi.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}
