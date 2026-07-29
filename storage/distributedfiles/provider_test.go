package distributedfiles_test

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/filesystem"
	"github.com/matthiasharzer/discord-drive/storage/distributedfiles"
	"github.com/matthiasharzer/discord-drive/util/fsutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFilesystemStorageProvider(maxSingleFileSize int64, directory string) storage.Provider {
	return distributedfiles.NewProvider(maxSingleFileSize, func(key string) chunk.Provider {
		return filesystem.NewProvider(path.Join(directory, key))
	})
}

func TestProvider_Save(t *testing.T) {
	t.Run("reads a saved file", func(t *testing.T) {
		directory, cleanup, err := fsutils.TemporaryDirectory()
		require.NoError(t, err)
		defer cleanup()

		provider := createFilesystemStorageProvider(int64(8), directory)

		err = provider.Save("testfile.txt", strings.NewReader("abcdefghijkl"))
		assert.NoError(t, err)

		reader, err := provider.Retrieve("testfile.txt")
		require.NoError(t, err)

		data, err := io.ReadAll(reader)
		require.NoError(t, err)

		assert.Equal(t, "abcdefghijkl", string(data))
	})

	t.Run("distributes content over multiple files", func(t *testing.T) {
		directory, cleanup, err := fsutils.TemporaryDirectory()
		require.NoError(t, err)
		defer cleanup()

		provider := createFilesystemStorageProvider(int64(8), directory)

		err = provider.Save("testfile.txt", strings.NewReader("abcdefghijkl"))
		assert.NoError(t, err)

		chunkFiles, err := os.ReadDir(path.Join(directory, "testfile.txt"))
		require.NoError(t, err)

		require.Equal(t, 2, len(chunkFiles))

		file1Bytes, err := os.ReadFile(path.Join(directory, "testfile.txt", chunkFiles[0].Name()))
		assert.NoError(t, err)
		file2Bytes, err := os.ReadFile(path.Join(directory, "testfile.txt", chunkFiles[1].Name()))
		assert.NoError(t, err)

		assert.Equal(t, "abcdefgh", string(file1Bytes))
		assert.Equal(t, "ijkl", string(file2Bytes))
	})
}
