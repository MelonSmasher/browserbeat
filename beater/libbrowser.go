package beater

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/elastic/beats/libbeat/logp"
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

// Utility function to determine if a string exists in an list
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Returns all Firefox profiles in a user's Firefox data path
func enumerateFirefoxProfiles(path string) []string {
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

// Returns Chrome history DB paths for each user on the machine
func getChromePaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string
		if isWindows() {
			userPath = filepath.Join("C:\\", "Users", user, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Google", "Chrome", "Default", "History")
		} else if isLinux() {
			userPath = filepath.Join("/home", user, ".config", "google-chrome", "Default", "History")
		}
		if stat, err := os.Stat(userPath); err == nil {
			if !stat.IsDir() {
				srcDestMap := new(srcAndDestPaths)
				srcDestMap.src = userPath
				srcDestMap.dest = filepath.Join(getScratchPath(user), "chrome.sqlite")

				ubhp := new(userBrowserHistoryPath)
				ubhp.user = user
				ubhp.paths = []srcAndDestPaths{*srcDestMap}
				paths = append(paths, *ubhp)
			}
		}
	}
	return paths
}

// Returns Firefox history DB paths for each user on the machine
func getFirefoxPaths(users []string) []userBrowserHistoryPath {
	var paths []userBrowserHistoryPath
	for _, user := range users {
		var userPath string

		if isWindows() {
			userPath = filepath.Join("C:", "Users", user, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
		} else if isMacos() {
			userPath = filepath.Join("/Users", user, "Library", "Application Support", "Firefox", "Profiles")
		} else if isLinux() {
			userPath = filepath.Join("/home", user, ".mozilla", "firefox")
		}
		if stat, err := os.Stat(userPath); err == nil {
			var fullPaths []srcAndDestPaths
			if stat.IsDir() {
				profiles := enumerateFirefoxProfiles(userPath)
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

// Returns all browser history paths for each user on the machine
func getBrowserHistoryPaths() systemBrowserHistoryPaths {
	users := enumerateUsers()
	histories := new(systemBrowserHistoryPaths)
	histories.chrome = getChromePaths(users)
	histories.firefox = getFirefoxPaths(users)
	histories.safari = getSafariPaths(users)
	return *histories
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

// Returns SQLite queries for each supported browser.
// Each query returns all data in the target DB
func getQueryMap() queryMap {
	qMap := new(queryMap)
	qMap.chrome = "SELECT datetime(last_visit_time/1000000-11644473600, 'unixepoch') as last_visited, url, title FROM urls;"
	qMap.firefox = "SELECT datetime(moz_historyvisits.visit_date/1000000,'unixepoch'), moz_places.url, moz_places.title FROM moz_places, moz_historyvisits WHERE moz_places.id = moz_historyvisits.place_id;"
	qMap.safari = "SELECT datetime(hv.visit_time + 978307200, 'unixepoch', 'localtime') as last_visited, hi.url, hv.title FROM history_visits hv, history_items hi WHERE hv.history_item = hi.id;"
	return *qMap
}

// Returns SQLite queries for each supported browser.
// Each query returns data since the supplied timestamp.
func getQueryMapSince(since string) queryMap {
	qMap := new(queryMap)
	qMap.chrome = "SELECT datetime(last_visit_time/1000000-11644473600, 'unixepoch') as last_visited, url, title FROM urls WHERE datetime(last_visit_time/1000000-11644473600, 'unixepoch') > datetime('" + since + "');"
	qMap.firefox = "SELECT datetime(moz_historyvisits.visit_date/1000000,'unixepoch'), moz_places.url, moz_places.title FROM moz_places, moz_historyvisits WHERE moz_places.id = moz_historyvisits.place_id AND datetime(moz_historyvisits.visit_date/1000000,'unixepoch') > datetime('" + since + "');"
	qMap.safari = "SELECT datetime(hv.visit_time + 978307200, 'unixepoch', 'localtime') as last_visited, hi.url, hv.title FROM history_visits hv, history_items hi WHERE hv.history_item = hi.id AND datetime(hv.visit_time + 978307200, 'unixepoch', 'localtime') > datetime('" + since + "');"
	return *qMap
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

// Performs a query on the target SQLite DB and returns the history entries
func getHistoryDataFromSqlite(file string, query string) []historyEntry {
	var entries []historyEntry
	db, err := sql.Open("sqlite3", file)
	checkErr(err)
	rows, err := db.Query(query)
	checkErr(err)
	for rows.Next() {
		var entry historyEntry
		err := rows.Scan(&entry.Date, &entry.Url, &entry.Title)
		checkErr(err)
		entry.UrlData = extractUrlData(entry.Url.String)
		entries = append(entries, entry)
	}
	err = rows.Err()
	checkErr(err)
	db.Close()
	return entries
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

// This function reads browser data for the target browser and returns it as an array
func readBrowserData(browsers systemBrowserHistoryPaths, browser string, hn hostnameObj, ipAddresses []string) []browserBeatData {
	// Init some vars
	var browserData []userBrowserHistoryPath
	var browserBeatDatas []browserBeatData
	dateFormat := "2006-01-02 15:04:05"

	// Chose the browser database paths based on the current browser
	if browser == "chrome" {
		browserData = browsers.chrome
	} else if browser == "firefox" {
		browserData = browsers.firefox
	} else if browser == "safari" {
		browserData = browsers.safari
	}

	if len(browserData) > 0 { // Are there any database => user pairs to go through?
		// Loop through each dataset in the collection
		for _, dataSet := range browserData {
			// Try to read the date of the last history entry from our state file
			stamp, err := readUserBrowserState(browser, dataSet.user)
			// Init the query vars
			query := ""
			qMap := getQueryMap()
			if err != nil { // We got an error when reading the state file
				// Set the SQL query to get all browser history
				qMap = getQueryMap()
				logp.Info("Sending all known " + browser + " history")
			} else { // If we did not get an error when reading the state file
				// Set the datetime stamp value
				stamp := string(stamp)
				// Change the SQL query strings to get history entries since the datetime stamp
				qMap = getQueryMapSince(stamp)
				logp.Info("Sending " + browser + " history since " + stamp)
			}
			// Based on our browser and the timestamp chose our SQL query
			if browser == "chrome" {
				query = qMap.chrome
			} else if browser == "firefox" {
				query = qMap.firefox
			} else if browser == "safari" {
				query = qMap.safari
			}

			// For each history database path object in the dataset
			for _, hist := range dataSet.paths {
				// Copy the history to our scratch directory.
				// We can't open the Sqlite DB if the browser has it locked.
				// So we make a copy of it in the OS's temp dir
				copyToScratch(hist.src, hist.dest, dataSet.user)
				// Read the history entries from the copied Sqlite DB using the query we picked earlier
				entries := getHistoryDataFromSqlite(hist.dest, query)
				// Go through each history entry
				for _, entry := range entries {
					// Convert the history SQL datetime to a native Golang datetime
					requestTime, err := time.Parse(dateFormat, entry.Date.String)
					if err != nil { // Fallback on the current datetime if we fail to convert
						requestTime = time.Now()
					}

					// Build the data object to be shipped by libbeat
					entryData := browserBeatData{
						requestTime,
						time.Now(),
						hn,
						event{
							dataset{
								urlInfo{
									entry.Date.String,
									entry.Url.String,
									entry.Title.String,
									entry.UrlData,
								},
								clientInfo{
									ipAddresses,
									dataSet.user,
									runtime.GOOS,
									browser,
									hn,
								},
							},
							"browserbeat-" + browser,
						},
					}

					// Append the structured history entry to the list that we return
					browserBeatDatas = append(browserBeatDatas, entryData)
				}
			}
			// If we got some data from the DB.
			// Store the datetime of that SQL entry for the next run so we can pick up where we left off.
			if len(browserBeatDatas) > 0 {
				lastIndex := len(browserBeatDatas) - 1
				lastBrowserEntry := browserBeatDatas[lastIndex]
				storeUserBrowserState(browser, dataSet.user, lastBrowserEntry.Event.Data.Entry.Date)
			}
		}
	}

	// Send the history entries array back
	return browserBeatDatas
}
