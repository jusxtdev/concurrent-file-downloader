// Responsible for managing the cli application
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLI struct {}

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

func NewCLIApplication() *CLI{
	return &CLI{}
}

func (c *CLI) Setup(){
	fmt.Println("Started CLI")
}

func (c *CLI) Loop() error {
	//TODO implement re-prompt on EOF and non-int values
	// idea - the numProcessed returned from fmt.Scanln return 0 if the input is either EOF or string, we can use that num value to create custom errors and use those to match in the Loop
	// but Scanln returns 1 for float values too and then raises and error
	for {
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
			fmt.Println(choice)
		case int(Exit):
			fmt.Println(choice)
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