package beater

import (
	"log"
	"os"
	"path/filepath"
)

func getSupportedBrowsers() []string {
	return []string{
		"chrome",
		"chrome-canary",
		"chrome-beta",
		"chrome-dev",
		"chromium",
		"edge",
		"edge-beta",
		"edge-dev",
		"edge-canary",
		"firefox",
		"safari",
		"vivaldi",
		"opera",
		"k-meleon",
		"brave",
	}
}

// Returns a list of chrome based browsers
// The following browsers have the same history db schema thus the same query is used
func getChromes() []string {
	return []string{
		"chrome",
		"chromium",
		"chrome-canary",
		"chrome-beta",
		"chrome-dev",
		"vivaldi",
		"opera",
		"brave",
		"edge",
		"edge-beta",
		"edge-dev",
		"edge-canary",
	}
}

// Returns a list of firefox based browsers
// The following browsers have the same history db schema thus the same query is used
func getFirefoxes() []string {
	return []string{
		"firefox",
		"k-meleon",
	}
}

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
