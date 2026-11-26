// neutral package to store url struct, to prevent cycle import errors
package url

type URL struct {
	Url string
	// file size in bytes
	TotalFileSize int64
	Filename string
	// mp4, png etc
	ContentType string			
	// "bytes" if download in parts is supported, "none" or "" otherwise 
	AcceptRanges string			
}