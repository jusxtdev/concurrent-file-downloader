// responsible for download management
// major concurrency stuff happens here
package download

import (
	"os"
	"sync"

	"github.com/jusxtdev/concurrent-file-download/internal/chunks"
	url_pkg "github.com/jusxtdev/concurrent-file-download/internal/url"
)

// function to setup the temporary download directory
func SetupDownloadDir() error {
	// check if the directory exists
	_, err := os.Stat("temp")
	if os.IsNotExist(err) {
		// create the directory
		err = os.Mkdir("temp", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}


// starts the chunk downloads, returns nil-error when done successfully
// workerCount is the same as no. of chunks
func StartDownload(chunks []chunks.Chunk, url url_pkg.URL) error{
	// create a temp folder for this file indside the temp folder created by setup
	err := createFolderForFile(url.Filename)
	if err != nil {
		return err
	}

	workerCount := len(chunks)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Go(func() {})
	}
	wg.Wait()
	return nil
}

func createFolderForFile(filename string) error {
	err := os.Mkdir("temp/"+filename, 0755)
	if err != nil {
		// TODO Handle duplicate files
		return err
	}
	return nil
}