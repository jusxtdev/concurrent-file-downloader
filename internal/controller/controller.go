// responsible for executing tasks after cmd/root menu choice is made
package controller

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Controller struct {}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) NewURL() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter URL : ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(line)
	return url, nil
}