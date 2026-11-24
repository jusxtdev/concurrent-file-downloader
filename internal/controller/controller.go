// responsible for executing tasks after cmd/root menu choice is made
package controller

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jusxtdev/concurrent-file-download/internal/helper"
	"github.com/jusxtdev/concurrent-file-download/internal/metadata"
)

type URL struct {
	Url string
	TotalFileSize float64
	Filename string
	ContentType string
	AcceptRanges string
}

type Controller struct {
	URL *URL
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

	// update the controller.URL
	c.URL = &URL{
		Url: url,
	} 

	// check metadata, whether URL has `Range request support`
	info, err := metadata.FetchMetadata(url)
	if err != nil {
		return err
	}
	fmt.Println(info)
	return nil
}