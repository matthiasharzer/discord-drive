package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/matthiasharzer/discord-drive/storage/chunk"
)

type ChunkProvider struct {
	storageDirectory string
}

func NewChunkProvider(storageDirectory string) chunk.Provider {
	return &ChunkProvider{
		storageDirectory: storageDirectory,
	}
}

func (p *ChunkProvider) getChunkPath(chunkIndex int) string {
	return filepath.Join(p.storageDirectory, strconv.Itoa(chunkIndex))
}

func (p *ChunkProvider) Writer(chunkIndex int) (io.WriteCloser, error) {
	chunkPath := p.getChunkPath(chunkIndex)

	err := os.MkdirAll(filepath.Dir(chunkPath), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
func (p *ChunkProvider) Reader(chunkIndex int) (io.ReadCloser, error) {
	chunkPath := p.getChunkPath(chunkIndex)

	file, err := os.OpenFile(chunkPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *ChunkProvider) ChunkExists(chunkIndex int) (bool, error) {
	chunkPath := p.getChunkPath(chunkIndex)
	_, err := os.Stat(chunkPath)
	return !os.IsNotExist(err), nil
}
