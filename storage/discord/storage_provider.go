package discord

import (
	"fmt"
	"io"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/discord"
	"github.com/matthiasharzer/discord-drive/storage/distributedfile"
)

const discordChunkSize = 19 * 1024 * 1024 // The limit is 20MB, but we leave some room for metadata and overhead

type storageProviderContext struct {
	client           *discordgo.Session
	storageChannelID string
	threads          map[string]string
	mu               *sync.RWMutex
}

func (c *storageProviderContext) fetchThreads() error {
	activeThread, err := c.client.ThreadsActive(c.storageChannelID)
	if err != nil {
		return fmt.Errorf("error fetching active threads: %w", err)
	}

	archivedThreads, err := c.client.ThreadsArchived(c.storageChannelID, nil, 100)
	if err != nil {
		return fmt.Errorf("error fetching archived threads: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, thread := range activeThread.Threads {
		c.threads[thread.Name] = thread.ID
	}
	for _, thread := range archivedThreads.Threads {
		c.threads[thread.Name] = thread.ID
	}

	return nil
}

func (c *storageProviderContext) createChunkProviderFunc() distributedfile.CreateChunkProviderFunc {
	return func(key string) (chunk.Provider, error) {
		c.mu.RLock()
		threadID, exists := c.threads[key]
		c.mu.RUnlock()

		if exists {
			return discord.NewProvider(c.client, threadID), nil
		}

		thread, err := c.client.ThreadStart(c.storageChannelID, key, discordgo.ChannelTypeGuildPublicThread, 60)
		if err != nil {
			return nil, fmt.Errorf("error creating thread for key %s: %w", key, err)
		}

		c.mu.Lock()
		defer c.mu.Unlock()
		c.threads[key] = thread.ID

		return discord.NewProvider(c.client, thread.ID), nil
	}
}

type StorageProvider struct {
	context                 *storageProviderContext
	distributedFileProvider storage.Provider
}

func NewStorageProvider(discordBotToken string, storageChannelID string) (storage.Provider, error) {
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", discordBotToken))
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	err = dg.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening Discord connection: %w", err)
	}

	context := &storageProviderContext{
		client:           dg,
		storageChannelID: storageChannelID,
		threads:          make(map[string]string),
		mu:               &sync.RWMutex{},
	}

	err = context.fetchThreads()
	if err != nil {
		return nil, fmt.Errorf("error fetching threads: %w", err)
	}

	return &StorageProvider{
		context:                 context,
		distributedFileProvider: distributedfile.NewStorageProvider(discordChunkSize, context.createChunkProviderFunc()),
	}, nil
}

func (s *StorageProvider) Close() error {
	err := s.distributedFileProvider.Close()
	if err != nil {
		return fmt.Errorf("error closing distributed file provider: %w", err)
	}

	err = s.context.client.Close()
	if err != nil {
		return fmt.Errorf("error closing Discord client: %w", err)
	}

	return nil
}

func (s *StorageProvider) Write(key string, data io.Reader) error {
	return s.distributedFileProvider.Write(key, data)
}
func (s *StorageProvider) Read(key string) (io.ReadCloser, error) {
	return s.distributedFileProvider.Read(key)
}
func (s *StorageProvider) Has(key string) (bool, error) {
	return s.distributedFileProvider.Has(key)
}
