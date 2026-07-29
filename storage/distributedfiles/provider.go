package distributedfiles

import (
	"fmt"
	"io"
	"sync"

	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
)

type distributeChunkReader struct {
	nextChunkIndex int
	chunkProvider  chunk.Provider
	currentReader  io.ReadCloser
}

func (r *distributeChunkReader) hasNextChunk() bool {
	return r.chunkProvider.ChunkExists(r.nextChunkIndex)
}

func (r *distributeChunkReader) advanceReader() error {
	if !r.hasNextChunk() {
		return nil
	}

	chunkFile, err := r.chunkProvider.Reader(r.nextChunkIndex)
	if err != nil {
		return err
	}
	r.currentReader = chunkFile
	r.nextChunkIndex++
	return nil
}

func (r *distributeChunkReader) Read(p []byte) (n int, err error) {
	if r.currentReader == nil {
		if !r.hasNextChunk() {
			return 0, io.EOF
		}
		err := r.advanceReader()
		if err != nil {
			return 0, err
		}
	}

	remainingBytes := len(p)
	for remainingBytes > 0 {
		chunkBytesRead, err := r.currentReader.Read(p[n:])
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

func (r *distributeChunkReader) Close() error {
	if r.currentReader == nil {
		return nil
	}
	return r.currentReader.Close()
}

type Provider struct {
	maxSingleFileSize   int64
	createChunkProvider func(key string) chunk.Provider
	locks               map[string]*sync.RWMutex
	mu                  *sync.RWMutex
}

func NewProvider(maxSingleFileSize int64, createChunkProvider func(key string) chunk.Provider) storage.Provider {
	return &Provider{
		maxSingleFileSize:   maxSingleFileSize,
		createChunkProvider: createChunkProvider,
		locks:               make(map[string]*sync.RWMutex),
		mu:                  &sync.RWMutex{},
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

func (p *Provider) Save(key string, data io.Reader) error {
	mu := p.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	chunkIndex := 0
	chunkProvider := p.createChunkProvider(key)

	b := make([]byte, p.maxSingleFileSize)
	for {
		n, err := data.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading data: %w", err)
		}
		writer, err := chunkProvider.Writer(chunkIndex)
		if err != nil {
			return fmt.Errorf("error getting writer for chunk %d: %w", chunkIndex, err)
		}
		_, err = writer.Write(b[:n])
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

	chunkProvider := p.createChunkProvider(key)

	return &distributeChunkReader{
		chunkProvider: chunkProvider,
	}, nil
}
