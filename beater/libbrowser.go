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

func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func isMacos() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

func isLinux() bool {
	if runtime.GOOS == "linux" {
		return true
	}
	return false
}

func mkDirP(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func getScratchPath(user string) string {
	scratchPath := filepath.Join(os.TempDir(), "tbBrowser", user)
	mkDirP(scratchPath)
	return scratchPath
}

func cleanScratchDir() {
	baseScratchDir := filepath.Join(os.TempDir(), "tbBrowser")
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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

func enumerateWindowsUsers() []string {
	return readUsers("C:\\Users")
}

func enumerateMacOSUsers() []string {
	return readUsers("/Users")
}

func enumerateLinuxUsers() []string {
	return readUsers("/home")
}

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

func getBrowserHistoryPaths() systemBrowserHistoryPaths {
	users := enumerateUsers()
	histories := new(systemBrowserHistoryPaths)
	histories.chrome = getChromePaths(users)
	histories.firefox = getFirefoxPaths(users)
	histories.safari = getSafariPaths(users)
	return *histories
}

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

func getQueryMap() queryMap {
	qMap := new(queryMap)
	qMap.chrome = "SELECT datetime(last_visit_time/1000000-11644473600, 'unixepoch') as last_visited, url, title FROM urls;"
	qMap.firefox = "SELECT datetime(moz_historyvisits.visit_date/1000000,'unixepoch'), moz_places.url, moz_places.title FROM moz_places, moz_historyvisits WHERE moz_places.id = moz_historyvisits.place_id;"
	qMap.safari = "SELECT datetime(hv.visit_time + 978307200, 'unixepoch', 'localtime') as last_visited, hi.url, hv.title FROM history_visits hv, history_items hi WHERE hv.history_item = hi.id;"
	return *qMap
}

func getQueryMapSince(since string) queryMap {
	qMap := new(queryMap)
	qMap.chrome = "SELECT datetime(last_visit_time/1000000-11644473600, 'unixepoch') as last_visited, url, title FROM urls WHERE datetime(last_visit_time/1000000-11644473600, 'unixepoch') > datetime('" + since + "');"
	qMap.firefox = "SELECT datetime(moz_historyvisits.visit_date/1000000,'unixepoch'), moz_places.url, moz_places.title FROM moz_places, moz_historyvisits WHERE moz_places.id = moz_historyvisits.place_id AND datetime(moz_historyvisits.visit_date/1000000,'unixepoch') > datetime('" + since + "');"
	qMap.safari = "SELECT datetime(hv.visit_time + 978307200, 'unixepoch', 'localtime') as last_visited, hi.url, hv.title FROM history_visits hv, history_items hi WHERE hv.history_item = hi.id AND datetime(hv.visit_time + 978307200, 'unixepoch', 'localtime') > datetime('" + since + "');"
	return *qMap
}

func copyHistory(src string, dst string, user string) (int64, error) {
	// Ensure that our scratch path exists
	getScratchPath(user)
	return copyFile(src, dst)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func extractUrlData(urlStr string) url.URL {
	urlData, err := url.Parse(urlStr)
	if err != nil {
		urlData = new(url.URL)
		return *urlData
	}
	return *urlData
}

func getHistoryData(file string, query string) []historyEntry {
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

func storeReadState(browser string, user string, stamp string) {
	mkDirP("states")
	data := []byte(stamp)
	err := ioutil.WriteFile(path.Join("states", browser+"_"+user+".state"), data, 0644)
	checkErr(err)
}

func readBrowserState(browser string, user string) ([]byte, error) {
	stamp, err := ioutil.ReadFile(path.Join("states", browser+"_"+user+".state"))
	return stamp, err
}

func readBrowserData(browsers systemBrowserHistoryPaths, browser string, hn hostnameObj, ipAddresses []string) []browserBeatData {
	var browserData []userBrowserHistoryPath
	if browser == "chrome" {
		browserData = browsers.chrome
	} else if browser == "firefox" {
		browserData = browsers.firefox
	} else if browser == "safari" {
		browserData = browsers.safari
	}
	var bites []browserBeatData
	dateFormat := "2006-01-02 15:04:05"

	if len(browserData) > 0 {
		for _, dataSet := range browserData {
			stamp, err := readBrowserState(browser, dataSet.user)
			query := ""
			qMap := getQueryMap()
			if err != nil {
				qMap = getQueryMap()
				logp.Info("Sending all known " + browser + " history")
			} else {
				stamp := string(stamp)
				qMap = getQueryMapSince(stamp)
				logp.Info("Sending " + browser + " history since " + stamp)
			}
			if browser == "chrome" {
				query = qMap.chrome
			} else if browser == "firefox" {
				query = qMap.firefox
			} else if browser == "safari" {
				query = qMap.safari
			}
			for _, hist := range dataSet.paths {
				copyHistory(hist.src, hist.dest, dataSet.user)
				entries := getHistoryData(hist.dest, query)
				for _, entry := range entries {
					requestTime, err := time.Parse(dateFormat, entry.Date.String)
					if err != nil {
						requestTime = time.Now()
					}

					bite := browserBeatData{
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

					bites = append(bites, bite)

				}
			}
			if len(bites) > 0 {
				last := len(bites) - 1
				lastBite := bites[last]
				storeReadState(browser, dataSet.user, lastBite.Event.Data.Entry.Date)
			}
		}
	}

	return bites
}
