package distributedfiles

import (
	"fmt"
	"io"
	"sync"

	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
)

type distributeChunkReader struct {
	nextChunkIndex        int
	chunkSize             int64
	currentChunkReadBytes int64
	chunkProvider         chunk.Provider
	currentReader         io.ReadCloser
}

func (r *distributeChunkReader) hasNextChunk() bool {
	return r.chunkProvider.ChunkExists(r.nextChunkIndex)
}

func (r *distributeChunkReader) advanceReader() error {
	if !r.hasNextChunk() {
		return nil
	}
	if r.currentReader != nil {
		err := r.currentReader.Close()
		if err != nil {
			return fmt.Errorf("error closing current reader: %w", err)
		}
	}

	chunkFile, err := r.chunkProvider.Reader(r.nextChunkIndex)
	if err != nil {
		return fmt.Errorf("error getting reader for chunk %d: %w", r.nextChunkIndex, err)
	}
	r.currentReader = chunkFile
	r.nextChunkIndex++
	r.currentChunkReadBytes = 0
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
		if err == io.EOF {
			return n, io.EOF
		}
		if err != nil {
			return n, fmt.Errorf("error reading chunk: %w", err)
		}
		remainingBytes -= chunkBytesRead
		r.currentChunkReadBytes += int64(chunkBytesRead)

		if r.currentChunkReadBytes >= r.chunkSize {
			err = r.advanceReader()
			if err != nil {
				return n, err
			}
		}

		wantsNextChunk := remainingBytes > 0
		if !wantsNextChunk {
			continue
		}
	}

	if !r.hasNextChunk() {
		r.currentReader = nil
		return n, nil
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
	chunkSize           int64
	createChunkProvider func(key string) chunk.Provider
	locks               map[string]*sync.RWMutex
	mu                  *sync.RWMutex
}

func NewProvider(chunkSize int64, createChunkProvider func(key string) chunk.Provider) storage.Provider {
	return &Provider{
		chunkSize:           chunkSize,
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

func (p *Provider) Write(key string, data io.Reader) error {
	mu := p.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	chunkIndex := 0
	chunkProvider := p.createChunkProvider(key)

	b := make([]byte, p.chunkSize)
	isEOF := false
	for !isEOF {
		n, err := data.Read(b)
		isEOF = err == io.EOF
		if err != nil && !isEOF {
			return fmt.Errorf("error reading data: %w", err)
		}
		if n == 0 {
			break
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

func (p *Provider) Read(key string) (io.ReadCloser, error) {
	mu := p.getMutex(key)
	mu.RLock()
	defer mu.RUnlock()

	chunkProvider := p.createChunkProvider(key)

	return &distributeChunkReader{
		chunkProvider: chunkProvider,
		chunkSize:     p.chunkSize,
	}, nil
}
