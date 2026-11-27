// responsible for fetching the HEAD metadata from the url
package metadata

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Info struct {				// Headers
	TotalFileSize int64			// Content-Length
	Filename string				// Content-Disposition
	AcceptRanges string			// Accept-Ranges
	ContentType string			// Content-Type
}

type DownloadType int
const (
	// request supports - Range request and Filesize is provided in the Content-Length
	ConcurrentDownload DownloadType = iota
	// request doesn't support - Range request and/or Filesize is not provided in the Content-Length
	SingleDownload
	// request supports - Range request but Filesize is not provided in the Content-Length
	NoFileSizeProvided
)

func FetchMetadata(url string) (Info, DownloadType, error) {
	info := Info{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return Info{}, 0, err
	}
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Info{}, 0, err
	}
	defer resp.Body.Close()

	//* REFER decisions.md -- what to do if HEAD request doesn't provide `TotalFileSize` 
	// Probe Request
	supportsRangeRequest, sizeFromProbe, err := ProbeRequest(url)
	if err != nil {
		return Info{}, 0, err
	}
	if !supportsRangeRequest {
		// the request doesn't support range requests, use single download
		return Info{}, SingleDownload, nil
	}

	sizeFromHead := resp.ContentLength
	switch {
	case sizeFromHead != -1:
		info.TotalFileSize = sizeFromHead
	case sizeFromProbe != -1:
		info.TotalFileSize = sizeFromProbe
	default:
		// filesize not available
		return Info{}, SingleDownload, nil
	}

	contentDispositionHeader := resp.Header.Get("Content-Disposition")
	info.Filename = extractFilenameFromContentDisposition(contentDispositionHeader)
	info.ContentType = resp.Header.Get("Content-Type")

	return info, ConcurrentDownload, nil
}

// makes a GET request with `Range` header and checks if the response code is 206 (Partial Content)
// Returns true if the request supports range requests, false otherwise
// also returns the size 
func ProbeRequest(url string) (bool, int64, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("Range", "bytes=0-0")

	resp, err := client.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return false, 0, nil
	}

	// range-request is supported, get the size from the response if the `Content-Range` header is present, or return -1 if not present
	reqRangeHeader := resp.Header.Get("Content-Range")
	if reqRangeHeader == "" {
		return true, -1, nil
	}
	var start, end, total int64
	_, err = fmt.Sscanf(reqRangeHeader, "bytes %d-%d/%d", &start, &end, &total)
	if err != nil {
		return false, 0, err
	}

	return true, total, nil
}

func GetFileSize(url string) (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	contentLength := resp.ContentLength
	if contentLength == -1 {
		return 0, fmt.Errorf("Content-Length header not found")
	}

	return resp.ContentLength, nil
}

func extractFilenameFromContentDisposition(header string) string {
	// process content length header to extract the filename
	// from header="attachment; filename", we want "filename"
	var filename string
	if strings.HasPrefix(header, "attachment;") {
		split := strings.Split(header, ";")
		// split[1] is "filename=<file1>"
		filenameQuoted := strings.TrimSpace(strings.Split(split[1], "=")[1])
		filename = strings.Trim(filenameQuoted, `"`)
	} else {
		filename = header
	}
	return filename
}