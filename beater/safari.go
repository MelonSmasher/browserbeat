package beater

import (
	"os"
	"path/filepath"
)

// Returns Safari history DB paths for each user on the machine
func getSafariPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Safari", "History.db")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "safari.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}
