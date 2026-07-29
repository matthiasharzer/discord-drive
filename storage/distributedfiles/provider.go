package distributedfiles

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/matthiasharzer/discord-drive/storage"
)

type distributedFileReader struct {
	chunkFiles     []string
	nextChunkIndex int
	chunkReader    io.ReadCloser
}

func (r *distributedFileReader) hasNextChunk() bool {
	return r.nextChunkIndex < len(r.chunkFiles)
}

func (r *distributedFileReader) advanceReader() error {
	if !r.hasNextChunk() {
		return nil
	}

	chunkPath := r.chunkFiles[r.nextChunkIndex]
	chunkFile, err := os.Open(chunkPath)
	if err != nil {
		return err
	}
	r.chunkReader = chunkFile
	r.nextChunkIndex++
	return nil
}

func (r *distributedFileReader) Read(p []byte) (n int, err error) {
	if r.chunkReader == nil {
		if r.nextChunkIndex >= len(r.chunkFiles) {
			return 0, io.EOF
		}
		err := r.advanceReader()
		if err != nil {
			return 0, err
		}
	}

	remainingBytes := len(p)
	for remainingBytes > 0 {
		chunkBytesRead, err := r.chunkReader.Read(p[n:])
		n += chunkBytesRead
		if err != nil && err != io.EOF {
			return n, fmt.Errorf("error reading chunk: %w", err)
		}
		remainingBytes -= chunkBytesRead
		if !r.hasNextChunk() {
			return n, io.EOF
		}
		err = r.advanceReader()
		if err != nil {
			return 0, err
		}
	}

	return n, nil
}

func (r *distributedFileReader) Close() error {
	if r.chunkReader == nil {
		return nil
	}
	return r.chunkReader.Close()
}

type Provider struct {
	maxSingleFileSize int64
	storageDirectory  string
	locks             map[string]*sync.RWMutex
	mu                *sync.RWMutex
}

func NewProvider(maxSingleFileSize int64, storageDirectory string) storage.Provider {
	return &Provider{
		maxSingleFileSize: maxSingleFileSize,
		storageDirectory:  storageDirectory,
		locks:             make(map[string]*sync.RWMutex),
		mu:                &sync.RWMutex{},
	}
}

func (p *Provider) getMutex(name string) *sync.RWMutex {
	p.mu.RLock()
	mu, ok := p.locks[name]
	p.mu.RUnlock()
	if !ok {
		mu = new(sync.RWMutex)
		p.mu.Lock()
		p.locks[name] = mu
		p.mu.Unlock()
	}
	return mu
}

func (p *Provider) getChunkPath(key string, chunkIndex int) string {
	return path.Join(p.storageDirectory, key, strconv.Itoa(chunkIndex))
}

func (p *Provider) writeChunk(key string, chunkIndex int, data []byte) error {
	chunkPath := p.getChunkPath(key, chunkIndex)

	err := os.MkdirAll(path.Dir(chunkPath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(chunkPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) Save(key string, data io.Reader) error {
	mu := p.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	chunkIndex := 0

	b := make([]byte, p.maxSingleFileSize)
	for {
		n, err := data.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading data: %w", err)
		}
		err = p.writeChunk(key, chunkIndex, b[:n])
		if err != nil {
			return fmt.Errorf("error writing chunk %d: %w", chunkIndex, err)
		}
		chunkIndex++
	}

	return nil
}

func (p *Provider) Retrieve(key string) (io.ReadCloser, error) {
	mu := p.getMutex(key)
	mu.RLock()
	defer mu.RUnlock()

	chunkFiles, err := os.ReadDir(path.Join(p.storageDirectory, key))
	if err != nil {
		return nil, err
	}

	var chunkFilePaths []string
	for _, dirEntry := range chunkFiles {
		chunkFilePath := path.Join(p.storageDirectory, key, dirEntry.Name())
		chunkFilePaths = append(chunkFilePaths, chunkFilePath)
	}

	return &distributedFileReader{
		chunkFiles: chunkFilePaths,
	}, nil
}
