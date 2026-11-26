// responsible for initializing chunks
package chunks

import (
	"fmt"

	url_pkg "github.com/jusxtdev/concurrent-file-download/internal/url"
)

var NoOfChunks = 4

type Chunk struct {
	ChunkId int
	StartByte int
	EndByte int
	Size int
}

func ChunkURL(urlInfo *url_pkg.URL) {
	baseSize := urlInfo.TotalFileSize / int64(NoOfChunks)
	// fmt.Println("Base size", baseSize)
	// fmt.Println("Total size", urlInfo.TotalFileSize)
	
	// create chunks
	var chunks []Chunk = make([]Chunk, NoOfChunks)
	for i := 0; i < NoOfChunks; i++ {
		start := i*int(baseSize)
		end := start + int(baseSize) - 1
		if end >= int(urlInfo.TotalFileSize) {
			end = int(urlInfo.TotalFileSize) - 1 
		}
		chunkSize := end - start + 1

		chunks[i] = Chunk{
			ChunkId: i,
			StartByte: start,
			EndByte: end,
			Size: chunkSize,
		}
	}

	for i := 0; i < NoOfChunks; i++ {
		fmt.Printf("%d : Start=%d -- End=%d -- Size=%d\n", chunks[i].ChunkId, chunks[i].StartByte, chunks[i].EndByte, chunks[i].Size)
	}
}