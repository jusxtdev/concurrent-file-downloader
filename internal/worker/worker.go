// worker responsible to download chunk
package worker

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jusxtdev/concurrent-file-download/internal/chunks"
	url_pkg "github.com/jusxtdev/concurrent-file-download/internal/url"
)

func DownloadChunk(chunk chunks.Chunk, url url_pkg.URL, chunkPath string, errChan chan error) error {
	fmt.Printf("Starting download for %d\n", chunk.ChunkId)
	fmt.Printf("- Start - %d | End - %d\n", chunk.StartByte, chunk.EndByte)
	// make get request
	req, err := http.NewRequest("GET", url.Url, nil)
	if err != nil {
		errChan <- err
	}

	// specify byte range
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", chunk.StartByte, chunk.EndByte))

	// recieve response
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		errChan <- err
	}
	defer resp.Body.Close()
	
	// if responded with something other than 206, then it is not a partial download
	if resp.StatusCode != http.StatusPartialContent {
		errChan <- fmt.Errorf("No partial Content - Status=%d", resp.StatusCode)
	}

	// save response in a temp chunk file
	file, err := os.Create(chunkPath)
	if err != nil {
		errChan <- err
	}
	fmt.Printf("%s file for %d\n", chunkPath, chunk.ChunkId)
	defer file.Close()
	bytesWritten, err := io.Copy(file, resp.Body)
	if err != nil {
		errChan <- err
	}
	if n := int64(chunk.StartByte - chunk.EndByte + 1); bytesWritten != n {
		errChan <- fmt.Errorf("expected bytes = %d, got = %d", n, bytesWritten)
	}
	fmt.Printf("completed download for %d\n", chunk.ChunkId)
	return nil
}