package main

import (
	"encoding/json"
	"github.com/charles7668/VideoDownloader/m3u8Downloader"
	"github.com/charles7668/VideoDownloader/m3u8Parser"
	"log"
	"os"
)

func main() {
	// get the command line arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Please provide a m3u8 url")
		return
	}

	// download m3u8 file
	m3u8Url := args[1]
	m3u8FilePath := "test.m3u8"
	log.Println("prepare to download m3u8 url : ", m3u8Url)
	log.Println("download m3u8 file path : ", m3u8FilePath)
	host, err := m3u8Downloader.DownloadM3U8File(m3u8Url, "test.m3u8")
	if err != nil {
		log.Fatal("download m3u8 file failed : " + err.Error())
		return
	}
	log.Println("download m3u8 file success")

	parseOutput := "temp_output_list.txt"
	log.Println("prepare to parse m3u8 file : ", m3u8FilePath)
	log.Println("parse output file path : ", parseOutput)
	// parse m3u8 file
	info, err := m3u8Parser.ParseM3U8File("test.m3u8", host.Relative, parseOutput)
	if err != nil {
		log.Fatal("parse failed : " + err.Error())
		return
	}
	log.Println("parse success")

	jsonString, err := json.Marshal(info)
	if err != nil {
		log.Fatal("json marshal failed : " + err.Error())
	}
	log.Printf("info: " + string(jsonString))
}
