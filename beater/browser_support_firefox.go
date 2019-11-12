package beater

import (
	"os"
	"path/filepath"
)

// Returns Firefox history DB paths for each user on the machine
func getFirefoxPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string

		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Firefox", "Profiles")
		} else if isLinux() {
			userPath = filepath.Join("/home", user, ".mozilla", "firefox")
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
							srcDestMap.dest = filepath.Join(getScratchPath(user), "firefox-"+profile+".sqlite")
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
