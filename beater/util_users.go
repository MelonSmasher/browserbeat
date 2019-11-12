package beater

import (
	"log"
	"os"
	"path/filepath"
)

// Reads user's from the home/users path on the filesystem
func readUsers(path string) []string {
	var dirNames []string
	winBuiltinUsers := []string{"Default", "Default User", "Public", "All Users"}
	macBuiltinUsers := []string{"Guest", "Shared"}
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
				if isWindows() {
					if !stringInSlice(name, winBuiltinUsers) {
						dirNames = append(dirNames, name)
					}
				} else if isMacos() {
					if !stringInSlice(name, macBuiltinUsers) {
						dirNames = append(dirNames, name)
					}
				} else {
					dirNames = append(dirNames, name)
				}
			}
		}
	}
	return dirNames
}

// Returns a list of Windows users
func enumerateWindowsUsers() []string {
	return readUsers("C:\\Users")
}

// Returns a list of macOS users
func enumerateMacOSUsers() []string {
	return readUsers("/Users")
}

// Returns a list of Linux users
func enumerateLinuxUsers() []string {
	return readUsers("/home")
}

// Returns a list of all users on the machine
func enumerateUsers() []string {
	if isWindows() {
		return enumerateWindowsUsers()
	} else if isMacos() {
		return enumerateMacOSUsers()
	} else if isLinux() {
		return enumerateLinuxUsers()
	} else {
		return []string{}
	}
}
