package beater

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
	"time"
)


type srcAndDestPaths struct {
	src  string
	dest string
}

type userBrowserHistoryPath struct {
	paths []srcAndDestPaths
	user  string
}

type systemBrowserHistoryPaths struct {
	chrome  []userBrowserHistoryPath
	firefox []userBrowserHistoryPath
	safari  []userBrowserHistoryPath
}

type queryMap struct {
	chrome  string
	firefox string
	safari  string
}

type hostnameObj struct {
	Hostname string `json:"hostname"`
	Short    string `json:"short"`
}

type historyEntry struct {
	Date    sql.NullString
	Url     sql.NullString
	Title   sql.NullString
	UrlData url.URL
}

type clientInfo struct {
	IpAddresses []string    `json:"ip_addresses"`
	User        string      `json:"user"`
	Platform    string      `json:"platform"`
	Browser     string      `json:"browser"`
	Hostname    hostnameObj `json:"Hostname"`
}

type biteData struct {
	Date    string  `json:"date"`
	Url     string  `json:"url"`
	Title   string  `json:"title"`
	UrlData url.URL `json:"url_data"`
}

type dataset struct {
	Entry  biteData   `json:"entry"`
	Client clientInfo `json:"client"`
}

type event struct {
	Data   dataset `json:"data"`
	Module string  `json:"module"`
}

type turkeyBite struct {
	Date      time.Time   `json:"@timestamp"`
	Processed time.Time   `json:"@processed"`
	Host      hostnameObj `json:"host"`
	Event     event       `json:"event"`
}
