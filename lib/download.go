package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// DownloadFromURL : Downloads the binary from the source url
func DownloadFromURL(installLocation string, url string) (string, error)  {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Download ", url)
	fileName = installLocation + fileName
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	return fileName, nil
}

