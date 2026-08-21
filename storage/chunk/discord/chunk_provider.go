package discord

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
)

type ChunkProvider struct {
	client          *discordgo.Session
	channelID       string
	messageIDLookup map[int]string
	latestMessageID string
	mu              *sync.RWMutex
}

func NewProvider(client *discordgo.Session, channelID string) chunk.Provider {
	return &ChunkProvider{
		client:          client,
		channelID:       channelID,
		messageIDLookup: make(map[int]string),
		mu:              &sync.RWMutex{},
	}
}

func (p *ChunkProvider) fetchMessages() error {
	numberOfMessages := 100
	hasMore := true
	p.mu.Lock()
	defer p.mu.Unlock()

	for hasMore {
		messages, err := p.client.ChannelMessages(p.channelID, numberOfMessages, p.latestMessageID, "", "")
		if err != nil {
			return fmt.Errorf("error fetching messages from Discord channel: %w", err)
		}
		hasMore = len(messages) == numberOfMessages
		if hasMore {
			p.latestMessageID = messages[len(messages)-1].ID
		}

		for _, message := range messages {
			if len(message.Attachments) != 1 {
				return fmt.Errorf("expected exactly one attachment per message, but got %d for message ID %s", len(message.Attachments), message.ID)
			}

			chunkIndex, err := strconv.Atoi(message.Attachments[0].Filename)
			if err != nil {
				return fmt.Errorf("error parsing chunk index from filename %s: %w", message.Attachments[0].Filename, err)
			}

			p.messageIDLookup[chunkIndex] = message.ID
		}
	}

	return nil
}

func (p *ChunkProvider) Writer(chunkIndex int) (io.WriteCloser, error) {
	chunkIndexStr := fmt.Sprintf("%d", chunkIndex)

	pr, pw := io.Pipe()

	go func() {
		defer func() {
			err := pr.Close()
			if err != nil {
				log.Printf("Error closing pipe reader: %v", err)
			}
		}()

		message, err := p.client.ChannelFileSend(p.channelID, chunkIndexStr, pr)
		if err != nil {
			log.Printf("Error sending file to Discord channel: %v", err)
		}

		p.mu.Lock()
		p.messageIDLookup[chunkIndex] = message.ID
		p.mu.Unlock()
	}()

	return pw, nil
}

func (p *ChunkProvider) Reader(chunkIndex int) (io.ReadCloser, error) {
	if len(p.messageIDLookup) == 0 {
		err := p.fetchMessages()
		if err != nil {
			return nil, fmt.Errorf("error fetching messages: %w", err)
		}
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	messageID, exists := p.messageIDLookup[chunkIndex]
	if !exists {
		return nil, fmt.Errorf("chunk index %d does not exist", chunkIndex)
	}

	message, err := p.client.ChannelMessage(p.channelID, messageID)
	if err != nil {
		return nil, fmt.Errorf("error fetching message with ID %s: %w", messageID, err)
	}

	if len(message.Attachments) != 1 {
		return nil, fmt.Errorf("expected exactly one attachment for message ID %s, but got %d", messageID, len(message.Attachments))
	}

	attachment := message.Attachments[0]
	resp, err := p.client.Client.Get(attachment.URL)
	if err != nil {
		return nil, fmt.Errorf("error downloading attachment from URL %s: %w", attachment.URL, err)
	}

	return resp.Body, nil
}

func (p *ChunkProvider) ChunkExists(chunkIndex int) (bool, error) {
	if len(p.messageIDLookup) == 0 {
		err := p.fetchMessages()
		if err != nil {
			return false, fmt.Errorf("error fetching messages: %w", err)
		}
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	_, exists := p.messageIDLookup[chunkIndex]
	return exists, nil
}

func (p *ChunkProvider) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}
