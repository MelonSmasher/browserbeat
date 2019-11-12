package beater

import (
	"log"
	"os"
	"path/filepath"
)

// Returns all Firefox profiles in a user's Firefox data path
func enumerateFFBrowserProfiles(path string) []string {
	var dirNames []string

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()
	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		info, err := os.Stat(filepath.Join(path, name))
		if err == nil {
			if info.IsDir() {
				dirNames = append(dirNames, name)
			}
		}
	}
	return dirNames
}

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
