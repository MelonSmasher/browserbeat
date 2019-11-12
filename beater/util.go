package beater

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Returns a structured hostname object that contains the full FQDN and the short hostname
func getHostname() hostnameObj {
	var hn hostnameObj
	name, err := os.Hostname()
	if err != nil {
		hn.Hostname = ""
		hn.Short = ""
	} else {
		hn.Hostname = name
		if strings.Contains(name, ".") {
			hn.Short = strings.Split(name, ".")[0]
		} else {
			hn.Short = name
		}
	}
	return hn
}

// Returns all IPs for this machine that are not associated with loopbacks
func getLocalIPs() []string {
	var localIps []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return localIps
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIps = append(localIps, ipnet.IP.String())
			}
		}
	}
	return localIps
}

// Lets us know if we're on Windows
func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

// Lets us know if we're on macOS
func isMacos() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

// Lets us know if we're on linux
func isLinux() bool {
	if runtime.GOOS == "linux" {
		return true
	}
	return false
}

// Utility function to determine if a string exists in an list
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Utility function to create a directory
func mkDirP(path string) {
	os.MkdirAll(path, os.ModePerm)
}

// Returns the path that should be used for working with the copied SQLite databases
func getScratchPath(user string) string {
	scratchPath := filepath.Join(os.TempDir(), "browserbeat", user)
	mkDirP(scratchPath)
	return scratchPath
}

// Removes and then creates a temp directory to copy each browser's database for each user
func cleanScratchDir() {
	baseScratchDir := filepath.Join(os.TempDir(), "browserbeat")
	src, err := os.Stat(baseScratchDir)
	if err == nil {
		if src.IsDir() {
			os.Remove(baseScratchDir)
		} else {
			os.Remove(baseScratchDir)
		}
	}
	mkDirP(baseScratchDir)
}

// Copies a file to it's destination
func copyFile(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// Copies a file to the OS's tmp dir nested under a folder named after the target user
func copyToScratch(src string, dst string, user string) (int64, error) {
	// Ensure that our scratch path exists
	getScratchPath(user)
	return copyFile(src, dst)
}

// Helper function to check errors and panic if needed
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Parses a URL and returns its components
func extractUrlData(urlStr string) url.URL {
	urlData, err := url.Parse(urlStr)
	if err != nil {
		urlData = new(url.URL)
		return *urlData
	}
	return *urlData
}

// Stores the timestamp of the last history entry from the current run for the target user => browser pair
func storeUserBrowserState(browser string, user string, stamp string) {
	mkDirP(path.Join("data", "states"))
	data := []byte(stamp)
	err := ioutil.WriteFile(path.Join("data", "states", browser+"_"+user+".state"), data, 0644)
	checkErr(err)
}

// Reads the timestamp of the last history entry from the last run for the target user => browser pair
func readUserBrowserState(browser string, user string) ([]byte, error) {
	stamp, err := ioutil.ReadFile(path.Join("data", "states", browser+"_"+user+".state"))
	return stamp, err
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
