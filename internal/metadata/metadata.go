package metadata

import (
	"net/http"
	"time"
)

type Info struct {				// Headers
	TotalFileSize int64			// Content-Length
	Filename string				// Content-Disposition
	AcceptRanges string			// Accept-Ranges
	ContentType string			// Content-Type
}
func FetchMetadata(url string) (Info, error) {
	info := Info{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return Info{}, err
	}
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Info{}, err
	}
	defer resp.Body.Close()

	info.TotalFileSize = int64(resp.ContentLength)
	info.AcceptRanges = resp.Header.Get("Accept-Ranges")
	info.Filename = resp.Header.Get("Content-Disposition")
	info.ContentType = resp.Header.Get("Content-Type")

	return info, err
}