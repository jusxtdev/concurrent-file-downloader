# Concurrent File Downloader

A command-line file downloader written in **Go** that downloads a single
remote file in parallel by splitting it into byte-range chunks,
downloading those chunks concurrently, and merging them into the
original file.

> **Project status:** MVP architecture and scope defined; implementation
> is in progress.

## Overview

Instead of downloading a file as one continuous response, the downloader
divides it into smaller byte ranges. A configurable worker pool
downloads those ranges concurrently. Each worker saves its chunk as a
separate temporary file. After every chunk has been downloaded and
validated successfully, the merger combines the chunk files in their
original order.

## How It Works

1.  **Input:** The user supplies a file URL, an output location, and
    download settings.
2.  **Metadata extraction:** The application retrieves the remote file's
    size and filename where available, then verifies that the server
    supports byte-range requests.
3.  **Chunk creation:** The chunk manager divides the file into
    contiguous, non-overlapping byte ranges.
4.  **Concurrent downloading:** A worker pool distributes chunk tasks to
    a fixed number of Go goroutines.
5.  **Temporary storage:** Each worker writes its chunk to its own
    temporary file.
6.  **Validation:** The application checks that each response contains
    the requested byte range and that the downloaded byte count matches
    the expected chunk size.
7.  **Merging:** Once all chunks succeed, the merger copies them
    sequentially in chunk-index order into a temporary output file.
8.  **Finalization:** The application verifies the final file size,
    moves it to the requested output filename, and removes temporary
    chunk files.

## Architecture

The application is organized around these components:

| Component | Responsibility |
| --- | --- |
| **Downloader Controller** | Coordinates the download lifecycle and decides when each stage can begin. |
| **Metadata Extractor** | Retrieves file metadata and verifies HTTP range support. |
| **Chunk Manager** | Calculates chunk boundaries and creates chunk tasks. |
| **Worker Pool** | Runs a configurable number of goroutines and distributes tasks through a channel. |
| **Chunk Downloader** | Sends HTTP range requests, validates responses, and writes chunk data. |
| **Chunk Storage** | Manages temporary chunk files and their lifecycle. |
| **Merger** | Combines validated chunks sequentially in their original order. |

### High-Level Flow

``` text
User Input
    |
    v
Downloader Controller
    |
    +--> Metadata Extractor
    |
    +--> Chunk Manager
    |
    v
Worker Pool
    |
    +--> Worker --> Chunk Downloader --> Temporary Chunk
    +--> Worker --> Chunk Downloader --> Temporary Chunk
    +--> Worker --> Chunk Downloader --> Temporary Chunk
    |
    v
Validate All Chunks
    |
    v
Sequential Merger
    |
    v
Final File
```

## Concurrency Model

The worker pool uses a fixed number of goroutines. Chunk tasks are
placed into a channel, and each worker takes one task at a time. When a
worker finishes a chunk, it can take another pending task.

The **number of chunks** and the **number of workers** are separate
settings. For example, a file may be divided into 20 chunks while only 4
workers run concurrently. Workers then process the chunks in batches.

A worker writes only to its own temporary chunk file. This keeps chunk
downloads independent and avoids multiple workers writing to the same
temporary file.

## Temporary File Strategy

Each download operation uses a dedicated temporary directory. A
conceptual layout is:

``` text
downloads/
├── .temp_<download>/
│   ├── chunk_000.part
│   ├── chunk_000
│   ├── chunk_001.part
│   ├── chunk_001
│   └── ...
└── <final-file>
```

A worker first writes to a partial file. After the response and byte
count have been validated and the file has been closed, the partial file
is renamed as a completed chunk. The merger only reads completed chunks.

The temporary output file is finalized as the requested output file only
after merging and final-size validation succeed.

## MVP Features

-   Download one file per operation.
-   Retrieve remote file metadata.
-   Verify HTTP byte-range support.
-   Split a file into configurable chunks.
-   Configure the number of concurrent workers.
-   Download chunks concurrently using goroutines.
-   Distribute tasks using Go channels.
-   Store chunks as separate temporary files.
-   Track aggregate download progress.
-   Validate HTTP status, content range, and downloaded byte count.
-   Merge chunks sequentially.
-   Verify the final file size.
-   Clean up temporary files after a successful merge.
-   Report errors and avoid merging when any chunk fails.

## Important Technical Considerations

### HTTP Range Requests

Chunk downloading depends on the server correctly supporting byte-range
requests. The downloader must verify this rather than relying only on an
`Accept-Ranges` response header.

For each chunk request, the application should:

-   Send a `Range` header with the requested inclusive byte boundaries.
-   Require `206 Partial Content`.
-   Validate the returned `Content-Range`.
-   Verify that the received byte count equals the expected chunk size.
-   Avoid treating a `200 OK` full-file response as a successful chunk
    response.

Some servers may not support range requests or may handle metadata
requests inconsistently. The MVP should report unsupported range
downloading clearly rather than silently producing an incorrect file.

### Chunk Boundaries

Chunk ranges must be contiguous, non-overlapping, and cover the entire
file. HTTP range end positions are inclusive. If the file size is not
evenly divisible by the chunk count, the final chunk may be smaller.

### Safe Merging

The merger must process chunks by their numeric index, not by an
arbitrary filesystem listing order. It should only run after every chunk
has been validated successfully.

### Compression and Byte Ranges

Byte-range offsets refer to the representation of the resource. To keep
offsets and downloaded bytes consistent, the downloader should request
identity encoding where possible and validate the server's response.

## Suggested Project Structure

The following is a proposed structure; it can evolve as implementation
progresses.

``` text
concurrent-file-downloader/
├── main.go
├── cmd/
│   └── root.go
├── internal/
│   ├── downloader/
│   │   └── downloader.go
│   ├── metadata/
│   │   └── metadata.go
│   ├── chunk/
│   │   └── chunk.go
│   ├── worker/
│   │   └── worker.go
│   ├── storage/
│   │   └── storage.go
│   └── merger/
│       └── merger.go
└── downloads/
```

## Development Roadmap

1.  Implement a basic single-file HTTP downloader.
2.  Retrieve file metadata and verify range support.
3.  Calculate chunk boundaries.
4.  Download one specified byte range to a temporary file.
5.  Implement the worker pool and concurrent chunk downloads.
6.  Add temporary chunk storage and validation.
7.  Implement sequential merging.
8.  Add aggregate progress reporting.
9.  Test error cases and verify the final output.

## Out of Scope for the MVP

The following are potential future enhancements, not MVP requirements:

-   Resume interrupted downloads.
-   Automatic retries for failed chunks.
-   Persistent download history.
-   Checksum verification against a trusted server-provided checksum.
-   Bandwidth throttling.
-   Downloading multiple files in one operation.
-   GUI or web interface.
-   Parallel writing directly into a shared output file.

## Concepts Practiced

This project is intended to provide practical experience with:

-   Go goroutines and channels.
-   Worker pool patterns and synchronization.
-   `sync.WaitGroup` and safe shared-state management.
-   HTTP requests, headers, status codes, and byte ranges.
-   Streaming data with `io.Reader` and `io.Writer`.
-   File creation, writing, closing, and renaming.
-   Error handling and resource cleanup.
-   Modular application design and testing.

## Building and Running

The command-line interface and its exact flags are not finalized yet.
Build and run instructions will be added once the CLI contract and
implementation are established.

## License

No license has been selected yet.
