package merge

import (
	"fmt"
	"os"
)

func MergeChunks(chunksPath string, outPath string) error {
	// create an out folder if not created
	_, err := os.Stat(outPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(outPath, 0755)
		if err != nil {
			return err
		}
	}
	// open each chunk files
	var tempFiles []string
	entries, err := os.ReadDir(chunksPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			tempFiles = append(tempFiles, entry.Name())
		}
	}

	fmt.Println("Your chunk files are -- ", tempFiles)
	// open / create the out file
	// read chunk files in order and append it's contents to outfile
	return nil
}