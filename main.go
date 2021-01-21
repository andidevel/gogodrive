package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/catfishlabs/gogodrive/config"
	"github.com/catfishlabs/gogodrive/gdrive"
	"google.golang.org/api/drive/v3"
)

const version = "0.0.1"

func getFileContentType(out *os.File) (string, error) {
	// From: https://golangcode.com/get-the-content-type-of-file/
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func evalDeletePattern(r *regexp.Regexp, deletePattern string) string {
	match := r.FindAllStringSubmatch(deletePattern, -1)
	today := time.Now()
	for i := 0; i < len(match); i++ {
		if match[i][1] == "today" {
			days, _ := strconv.Atoi(match[i][2])
			modfiedTime := today.AddDate(0, 0, days)
			deletePattern = strings.Replace(deletePattern, match[i][0], modfiedTime.Format(time.RFC3339), -1)
		}
	}
	return deletePattern
}

func main() {
	fmt.Printf("GoGoDrive - v%s\n", version)
	var (
		confFile    string
		uploadFile  string
		outFileName string
	)
	flag.StringVar(&confFile, "c", "configuration.json", "JSON configuration file path")
	flag.StringVar(&uploadFile, "i", "", "A file to upload")
	flag.StringVar(&outFileName, "o", "", "Output file name")
	flag.Parse()

	conf, err := config.ReadConfigFile(confFile)
	if err != nil {
		panic(err)
	}
	service, err := gdrive.GetService(conf.CredentialFile, conf.TokenFile)
	if err != nil {
		panic(err)
	}
	// Delete pattern regexp
	deletePatternRegexp, _ := regexp.Compile("{(.*)([+-][0-9]+)}")

	// First apply delete pattern, if any
	if conf.DeletePattern != "" {
		deletePatternStr := evalDeletePattern(deletePatternRegexp, conf.DeletePattern)
		fmt.Printf("Delete Pattern: %s\n", deletePatternStr)
		fileList, err := gdrive.SearchFile(service, deletePatternStr)
		if err == nil {
			if len(fileList.Files) > 0 {
				fmt.Printf("Found %d file(s).\n", len(fileList.Files))
				for _, i := range fileList.Files {
					fmt.Printf("Deleting %s\n", i.Name)
					err = gdrive.DeleteFile(service, i.Id)
				}
			} else {
				fmt.Println(" -> No files found!")
			}
		} else {
			fmt.Printf(" -> Error: %v\n", err)
		}
	}
	// Now, upload the new file
	f, err := os.Open(uploadFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	mimeType, err := getFileContentType(f)
	if err != nil {
		panic(err)
	}

	var uploadFilename string
	if outFileName != "" {
		uploadFilename = outFileName
	} else {
		uploadFilename = path.Base(uploadFile)
	}
	fmt.Printf("Uploading %s -> %s [%s]...\n", uploadFile, uploadFilename, mimeType)

	var file *drive.File
	for i := 0; i < conf.UploadMaxTries; i++ {
		fmt.Printf(" - Attempt %d/%d...\n", i+1, conf.UploadMaxTries)
		f.Seek(0, 0)
		file, err = gdrive.CreateFile(service, uploadFilename, mimeType, f, "root")
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Printf("-> Error: %v\n", err)
	} else {
		fmt.Printf("Successfully uploaded: %s.\n", file.Name)
	}
	fmt.Println("Done.")
}
