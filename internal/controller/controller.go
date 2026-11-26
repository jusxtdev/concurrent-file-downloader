// responsible for executing tasks after cmd/root menu choice is made
package controller

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jusxtdev/concurrent-file-download/internal/chunks"
	"github.com/jusxtdev/concurrent-file-download/internal/helper"
	"github.com/jusxtdev/concurrent-file-download/internal/metadata"
	url_pkg "github.com/jusxtdev/concurrent-file-download/internal/url"
)


type Controller struct {
	URL *url_pkg.URL
}

var (
	ErrInvalidURL = fmt.Errorf("invalid URL")
)

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) NewURL() error {
	// get the url from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter URL : ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	url := strings.TrimSpace(line)

	// check if URL is a valid HTTP or HTTPS url
	valid := helper.IsValidURL(url)
	if !valid {
		return ErrInvalidURL
	}

	// check metadata, whether URL has `Range request support`
	info, downloadType, err := metadata.FetchMetadata(url)
	if err != nil {
		return err
	}

	switch downloadType {
	case metadata.ConcurrentDownload:
		fmt.Println("The URL supports Range request, will use concurrent download")
	case metadata.SingleDownload:
		fmt.Println("The URL doesn't support Range request, will use single download")
	case metadata.NoFileSizeProvided:
		fmt.Println("The URL supports Range request but doesn't provide TotalFileSize, will use single download")
	}

	// update the controller.URL
	c.URL = &url_pkg.URL{
		Url: url,
		TotalFileSize: info.TotalFileSize,
		Filename: info.Filename,
		ContentType: info.ContentType,
		AcceptRanges: info.AcceptRanges,
	} 
	// fmt.Println(c.URL.TotalFileSize)
	chunks.ChunkURL(c.URL)
	return nil
}