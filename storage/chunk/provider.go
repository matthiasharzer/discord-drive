package chunk

import "io"

type Provider interface {
	Writer(chunkIndex int) (io.WriteCloser, error)
	Reader(chunkIndex int) (io.ReadCloser, error)
	ChunkExists(chunkIndex int) bool
}
