// Responsible for managing the cli application
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jusxtdev/concurrent-file-download/internal/controller"
)

type CLI struct {
	controller *controller.Controller
}

type Option int

const (
	NewURL Option = iota + 1
	Exit
)
var choices = map[Option]string{
	NewURL: "Enter URL to download",
	Exit: "Exit",
}

var ErrInvalidChoice = fmt.Errorf("Invalid choice")

func NewCLIApplication(c *controller.Controller) *CLI{
	return &CLI{
		controller: c,
	}
}

func (cli *CLI) Setup(){
}

// CLI application loop
func (cli *CLI) Loop() error {
	for {
		// print menu and take int input for choice
		// reprompt on invalid int choice or string input
		// EOF ignored for now
		choice, err := menu()
		if err != nil {
			if errors.Is(err, ErrInvalidChoice){
				// invalid choice entered, set choice to -1 so it goes to default case i.e. ErrInvalidChoice
				choice = -1
			} else {
				return err
			}
		}

		switch choice {
		case int(NewURL):
			err := cli.controller.NewURL()
			if err != nil {
				if errors.Is(err, controller.ErrInvalidURL){
					fmt.Printf("ERR : %s\n", err.Error())
				} else {
					return err
				}
			}
		case int(Exit):
			return nil
		default:
			fmt.Println(ErrInvalidChoice.Error())
		}
	}
}

/* = = = = = = = = = = = = = = = = = = = */
func menu() (int, error) {
	for option, choice := range choices{
		fmt.Printf("%d : %s\n", option, choice)
	}

	fmt.Printf("Enter choice: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	choice, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return 0, ErrInvalidChoice
	}
	return choice, nil
}