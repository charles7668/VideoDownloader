package m3u8Downloader

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type UrlInfo struct {
	Host     string
	Relative string
}

// closeFile closes the file
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

// getUrlHost returns the host of the url
func getUrlHost(videoUrl string) string {
	u, err := url.Parse(videoUrl)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return u.Scheme + "://" + u.Host
}

// DownloadM3U8File downloads the m3u8 file and returns the host of the m3u8 file
// url: the url of the m3u8 file
// filePath: the file path to save the m3u8 file
func DownloadM3U8File(url string, outputFile string) (UrlInfo, error) {
	host := getUrlHost(url)
	info := UrlInfo{}
	if host == "" {
		return info, errors.New("not a valid url")
	}
	extension := filepath.Ext(url)
	if extension != ".m3u8" {
		return info, errors.New("not a m3u8 file")
	}
	resp, err := http.Get(url)

	if err != nil {
		return info, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/vnd.apple.mpegurl" {
		return info, errors.New("not a HLS protocal file")
	}
	file, err := os.Create(outputFile)
	if err != nil {
		return info, err
	}
	// don't forget to close the file
	defer closeFile(file)

	writtenSize, err := io.Copy(file, resp.Body)
	log.Printf("writtenSize: %d", writtenSize)

	lastSlashIndex := strings.LastIndex(url, "/")
	info.Relative = url[:lastSlashIndex]
	info.Host = host

	return info, nil
}
