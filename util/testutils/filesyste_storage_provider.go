package testutils

import (
	"path"
	"testing"

	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/filesystem"
	"github.com/matthiasharzer/discord-drive/storage/distributedfiles"
	"github.com/matthiasharzer/discord-drive/util/fsutils"
	"github.com/stretchr/testify/require"
)

func FilesystemStorageProvider(t *testing.T, maxSingleFileSize int64, directory string) storage.Provider {
	t.Helper()
	return distributedfiles.NewProvider(maxSingleFileSize, func(key string) chunk.Provider {
		return filesystem.NewProvider(path.Join(directory, key))
	})
}

func TempDirFilesystemStorageProvider(t *testing.T, maxSingleFileSize int64) (storage.Provider, func()) {
	directory, cleanup, err := fsutils.TemporaryDirectory()
	require.NoError(t, err)

	return FilesystemStorageProvider(t, maxSingleFileSize, directory), cleanup
}
