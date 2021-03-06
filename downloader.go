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
	// ImmuniURL is the base url for the Immuni app
	ImmuniURL = "https://get.immuni.gov.it"
	// SwissCovidURL is the base url for the SwissCovid app
	SwissCovidURL = "https://www.pt.bfs.admin.ch"
)

// DownloaderFactory returns the Downloader for the specified app
func DownloaderFactory(app string) (Downloader, error) {
	switch app {
	case "immuni":
		return ImmuniDownloader{}, nil
	case "swisscovid":
		return SwissCovidDownloader{}, nil
	}
	return nil, fmt.Errorf("unknown app [%s]", app)
}

// ImmuniDownloader is the downloader for the Immuni app
type ImmuniDownloader struct{}

// GetLatestExport returns the latest Immuni export
func (d ImmuniDownloader) GetLatestExport() (string, error) {
	resp, err := http.Get(ImmuniURL + "/v1/keys/index")
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

// GetURL returns the Immuni URL where to download the export
func (d ImmuniDownloader) GetURL(export string) string {
	return ImmuniURL + "/v1/keys/" + export
}

// SwissCovidDownloader is the downloader for the Immuni app
type SwissCovidDownloader struct{}

// GetLatestExport returns the latest SwissCovid export
func (d SwissCovidDownloader) GetLatestExport() (string, error) {
	retry := 0

	for retry < 3 {
		now := time.Now()
		nowMidnight := time.Date(now.Year(), now.Month(), now.Day()-retry, 0, 0, 0, 0, time.UTC)
		latestExport := strconv.Itoa(int(nowMidnight.Unix() * 1000))

		url := d.GetURL(latestExport)
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return latestExport, nil
		}
		retry++
	}

	return "", fmt.Errorf("error getting latest swisscovid export")
}

// GetURL returns the SwissCovid URL where to download the export
func (d SwissCovidDownloader) GetURL(export string) string {
	return SwissCovidURL + "/v1/gaen/exposed/" + export
}
