package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader interface {
	GetLatestExport() (string, error)
	GetURL(export string) string
}

func Download(workDir, app string) error {
	dwln, err := DownloaderFactory(app)
	if err != nil {
		return err
	}

	path := filepath.Join(workDir, app)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	latest, _ := dwln.GetLatestExport()
	latestURL := dwln.GetURL(latest)

	filepath := filepath.Join(workDir, app, latest)
	filepathZip := filepath + ".zip"

	if err = DownloadZip(latestURL, filepathZip); err != nil {
		return err
	}

	if _, err := Unzip(filepathZip, filepath); err != nil {
		return err
	}

	return os.Remove(filepathZip)
}

func DownloadZip(url, zipPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading zip: status code %d", resp.StatusCode)
	}

	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
