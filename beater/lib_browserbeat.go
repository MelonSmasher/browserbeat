package beater

import (
	"runtime"
	"time"
)

// Returns all browser history paths for each user on the machine
func getBrowserHistoryPaths() systemBrowserHistoryPaths {
	users := enumerateUsers()
	histories := new(systemBrowserHistoryPaths)
	histories.chrome = getChromePaths(users)
	histories.chromium = getChromiumPaths(users)
	histories.chromeCanary = getChromeCanaryPaths(users)
	histories.chromeBeta = getChromeBetaPaths(users)
	histories.chromeDev = getChromeDevPaths(users)
	histories.firefox = getFirefoxPaths(users)
	histories.safari = getSafariPaths(users)
	histories.vivaldi = getVivaldiPaths(users)
	histories.opera = getOperaPaths(users)
	histories.kmeleon = getKmeleonPaths(users)
	histories.brave = getBravePaths(users)
	histories.edge = getEdgePaths(users)
	histories.edgeBeta = getEdgeBetaPaths(users)
	histories.edgeDev = getEdgeDevPaths(users)
	histories.edgeCanary = getEdgeCanaryPaths(users)
	return *histories
}

// Returns the browser data paths for this machine and target browser combo
func chooseBrowserDataPath(browser string, browsers systemBrowserHistoryPaths) []userBrowserHistoryPath {
	var none []userBrowserHistoryPath
	switch browser {
	case "chrome":
		return browsers.chrome
	case "firefox":
		return browsers.firefox
	case "safari":
		return browsers.safari
	case "chromium":
		return browsers.chromium
	case "chrome-canary":
		return browsers.chromeCanary
	case "chrome-beta":
		return browsers.chromeBeta
	case "chrome-dev":
		return browsers.chromeDev
	case "vivaldi":
		return browsers.vivaldi
	case "opera":
		return browsers.opera
	case "k-meleon":
		return browsers.kmeleon
	case "brave":
		return browsers.brave
	case "edge":
		return browsers.edge
	case "edge-canary":
		return browsers.edgeCanary
	case "edge-beta":
		return browsers.edgeBeta
	case "edge-dev":
		return browsers.edgeDev
	default:
		return none
	}
}

// This function reads browser data for the target browser and returns it as an array
func readBrowserData(browsers systemBrowserHistoryPaths, browser string, hn hostnameObj, ipAddresses []string) []browserBeatData {
	// Init some vars
	var browserData []userBrowserHistoryPath
	var browserBeatDatas []browserBeatData
	dateFormat := "2006-01-02 15:04:05"
	// Chrome based browsers
	chromes := getChromes()
	// Firefox based browsers
	firefoxes := getFirefoxes()

	// Chose the browser database paths based on the current browser
	browserData = chooseBrowserDataPath(browser, browsers)

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
			} else { // If we did not get an error when reading the state file
				// Set the datetime stamp value
				stamp := string(stamp)
				// Change the SQL query strings to get history entries since the datetime stamp
				qMap = getQueryMapSince(stamp)
			}

			// Based on our browser and the timestamp chose our SQL query
			if stringInSlice(browser, chromes) {
				query = qMap.chrome
			} else if stringInSlice(browser, firefoxes) {
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
