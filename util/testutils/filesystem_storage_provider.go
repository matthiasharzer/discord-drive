package testutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/filesystem"
	"github.com/matthiasharzer/discord-drive/storage/distributedfile"
	"github.com/matthiasharzer/discord-drive/util/fsutils"
)

func FilesystemStorageProvider(t *testing.T, chunkSize int64, directory string) storage.Provider {
	t.Helper()
	return distributedfile.NewStorageProvider(chunkSize, func(key string) (chunk.Provider, error) {
		return filesystem.NewChunkProvider(filepath.Join(directory, key)), nil
	})
}

func TempDirFilesystemStorageProvider(t *testing.T, chunkSize int64) (storage.Provider, func()) {
	directory, cleanup, err := fsutils.TemporaryDirectory()
	require.NoError(t, err)

	return FilesystemStorageProvider(t, chunkSize, directory), cleanup
}
