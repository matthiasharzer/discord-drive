package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/matthiasharzer/discord-drive/storage/chunk"
)

type Provider struct {
	storageDirectory string
}

func NewProvider(storageDirectory string) chunk.Provider {
	return &Provider{
		storageDirectory: storageDirectory,
	}
}

func (p *Provider) getChunkPath(chunkIndex int) string {
	return filepath.Join(p.storageDirectory, strconv.Itoa(chunkIndex))
}

func (p *Provider) Writer(chunkIndex int) (io.WriteCloser, error) {
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
func (p *Provider) Reader(chunkIndex int) (io.ReadCloser, error) {
	chunkPath := p.getChunkPath(chunkIndex)

	file, err := os.OpenFile(chunkPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *Provider) ChunkExists(chunkIndex int) bool {
	chunkPath := p.getChunkPath(chunkIndex)
	_, err := os.Stat(chunkPath)
	return !os.IsNotExist(err)
}
