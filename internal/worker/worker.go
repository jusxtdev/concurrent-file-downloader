// worker responsible to download chunk
package worker

import (
	"fmt"

	"github.com/jusxtdev/concurrent-file-download/internal/chunks"
)

func DownloadChunk(chunk chunks.Chunk) {
	fmt.Printf("Downloaded %d\n", chunk.ChunkId)
}