package m3u8Parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type M3U8Info struct {
	Host               string
	FilePath           string
	TsListFileLocation string
	Version            string
	PlayListType       string
	MediaSequence      string
	TargetDuration     string
	TsFileCount        uint
}

// parseKeyWord parses the m3u8 file and returns the M3U8Info struct
func parseKeyWord(info *M3U8Info, line string) {
	if strings.HasPrefix(line, "#EXT-X-VERSION:") {
		info.Version = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
	} else if strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:") {
		info.PlayListType = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
	} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
		info.MediaSequence = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
	} else if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
		info.TargetDuration = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
	}
}

// closeFile closes the file , generally used for defer
// file: the file to be closed
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// ParseM3U8File parses the m3u8 file and returns the M3U8Info struct
// fileName: the m3u8 file path
// host: the host of the m3u8 file
// outputFilePath: the output file path of the ts list
func ParseM3U8File(fileName string, host string, outputFilePath string) (*M3U8Info, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Fatal(err)
		return nil, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer closeFile(file)

	info := M3U8Info{}
	info.TsListFileLocation = outputFilePath
	info.Host = host

	scanner := bufio.NewScanner(file)
	tsList := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			parseKeyWord(&info, line)
		} else {
			if strings.HasPrefix(line, "http") {
				tsList = append(tsList, line)
			} else {
				tsList = append(tsList, host+"/"+line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, nil
	}

	info.TsFileCount = uint(len(tsList))

	outFile, err := os.Create(outputFilePath)
	defer closeFile(outFile)
	writer := bufio.NewWriter(outFile)
	for _, line := range tsList {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	err = writer.Flush()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, nil
	}
	return &info, nil
}
