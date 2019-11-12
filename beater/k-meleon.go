package beater

import (
	"os"
	"path/filepath"
)

// Returns K-Meleon history DB paths for each user on the machine
func getKmeleonPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string

		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Roaming", "K-Meleon")
		}
		if stat, err := os.Stat(userPath); err == nil {
			var fullPaths []srcAndDestPaths
			if stat.IsDir() {
				profiles := enumerateFFBrowserProfiles(userPath)
				for _, profile := range profiles {
					p := filepath.Join(userPath, profile, "places.sqlite")
					if pStat, err := os.Stat(p); err == nil {
						if !pStat.IsDir() {
							srcDestMap := new(srcAndDestPaths)
							srcDestMap.src = p
							srcDestMap.dest = filepath.Join(getScratchPath(user), "k-meleon-"+profile+".sqlite")
							fullPaths = append(fullPaths, *srcDestMap)
						}
					}
				}
			}
			ubhp := new(userBrowserHistoryPath)
			ubhp.user = user
			ubhp.paths = fullPaths
			paths = append(paths, *ubhp)
		}
	}
	return paths
}
