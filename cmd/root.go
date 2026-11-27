// Responsible for managing the cli application
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/jusxtdev/concurrent-file-download/internal/controller"
	"github.com/jusxtdev/concurrent-file-download/internal/download"
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

func (cli *CLI)ShutDown() error {
	err := cleanupTemp()
	if err != nil {
		return err
	}
	return nil
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

func cleanupTemp() error{
	// get list of dirs in temp
	entries, err := os.ReadDir(download.TempFolderBase)
	if err != nil {
		return err
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir(){
			dirs = append(dirs, entry.Name())
		}
	}

	// check which sub-dir is empty
	var dirsToDelete []string
	for _, subdir := range dirs {
		dirpath := fmt.Sprintf("%s/%s", download.TempFolderBase, subdir)
		empty, err := isDirEmpty(dirpath)
		if err != nil {
			return err
		}
		if empty {
			dirsToDelete = append(dirsToDelete, dirpath)
		}
	}
	// delete that dir 
	for _, dirPath := range dirsToDelete {
		err := os.RemoveAll(dirPath)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, nil
}