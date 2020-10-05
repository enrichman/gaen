package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	IMMUNI_URL      = "https://get.immuni.gov.it"
	SWISS_COVID_URL = "https://www.pt.bfs.admin.ch"
)

func DownloaderFactory(app string) (Downloader, error) {
	switch app {
	case "immuni":
		return ImmuniDownloader{}, nil
	case "swisscovid":
		return SwissCovidDownloader{}, nil
	}
	return nil, fmt.Errorf("unknown app [%s]", app)
}

type ImmuniDownloader struct{}

func (d ImmuniDownloader) GetLatestExport() (string, error) {
	resp, err := http.Get(IMMUNI_URL + "/v1/keys/index")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting latest immuni export. Status code %d", resp.StatusCode)
	}

	var m map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}

	latest, ok := m["newest"]
	if !ok {
		return "", errors.New("unknown 'newest' field")
	}

	return strconv.Itoa(latest), nil
}

func (d ImmuniDownloader) GetURL(export string) string {
	return IMMUNI_URL + "/v1/keys/" + export
}

type SwissCovidDownloader struct{}

func (d SwissCovidDownloader) GetLatestExport() (string, error) {
	now := time.Now()
	nowMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return strconv.Itoa(int(nowMidnight.Unix() * 1000)), nil
}

func (d SwissCovidDownloader) GetURL(export string) string {
	return SWISS_COVID_URL + "/v1/gaen/exposed/" + export
}
