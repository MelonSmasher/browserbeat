package beater

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

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
