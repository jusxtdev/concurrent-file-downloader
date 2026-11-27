// responsible for download management
// major concurrency stuff happens here
package download

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jusxtdev/concurrent-file-download/internal/chunks"
	url_pkg "github.com/jusxtdev/concurrent-file-download/internal/url"
	"github.com/jusxtdev/concurrent-file-download/internal/worker"
)

var TempFolderBase = "temp"

// function to setup the temporary download directory
func SetupDownloadDir() error {
	// check if the directory exists
	_, err := os.Stat(TempFolderBase)
	if os.IsNotExist(err) {
		// create the directory
		err = os.Mkdir(TempFolderBase, 0755)
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

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	for i := 0; i < workerCount; i++ {
		chunk := chunks[i]
		wg.Go(func() {
			worker.DownloadChunk(chunk)
		})
	}
	wg.Wait()
	return nil
}

func createFolderForFile(filename string) error {
	// Format: YYYYMMDD_HHMMSS (e.g., 20261127_194700)
	timestamp := time.Now().Format("20060102_150405")
	folderName := fmt.Sprintf("%s/%s_%s", TempFolderBase, filename, timestamp)

	err := os.Mkdir(folderName, 0755)
	if err != nil {
		return err
	}
	return nil
}