package distributedfile_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matthiasharzer/discord-drive/util/fsutils"
	"github.com/matthiasharzer/discord-drive/util/testutils"
)

func TestProvider_Save(t *testing.T) {
	t.Run("reads a clean divided file", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(4))
		defer cleanup()

		err := provider.Write("testfile.txt", strings.NewReader("abcdefghijkl")) // 12 byte = 3 chunks
		assert.NoError(t, err)

		reader, err := provider.Read("testfile.txt")
		require.NoError(t, err)
		defer func() {
			err := reader.Close()
			require.NoError(t, err)
		}()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)

		assert.Equal(t, "abcdefghijkl", string(data))
	})

	t.Run("reads a clean divided file with an almost full chunk", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(4))
		defer cleanup()

		err := provider.Write("testfile.txt", strings.NewReader("abcdefghijk")) // 11 byte = 3 chunks; last chunk has only 3 bytes
		assert.NoError(t, err)

		reader, err := provider.Read("testfile.txt")
		require.NoError(t, err)
		defer func() {
			err := reader.Close()
			require.NoError(t, err)
		}()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)

		assert.Equal(t, "abcdefghijk", string(data))
	})

	t.Run("reads a clean divided file with an almost empty chunk", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(4))
		defer cleanup()

		err := provider.Write("testfile.txt", strings.NewReader("abcdefghijklm")) // 13 byte = 4 chunks; last chunk has only 1 byte
		assert.NoError(t, err)

		reader, err := provider.Read("testfile.txt")
		require.NoError(t, err)
		defer func() {
			err := reader.Close()
			require.NoError(t, err)
		}()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)

		assert.Equal(t, "abcdefghijklm", string(data))
	})

	t.Run("reads large saved file", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(8))
		defer cleanup()

		content := `
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore
magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd
gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing
elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos
et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor
sit amet.
`

		err := provider.Write("testfile.txt", strings.NewReader(content))
		require.NoError(t, err)

		reader, err := provider.Read("testfile.txt")
		require.NoError(t, err)
		defer func() {
			err := reader.Close()
			require.NoError(t, err)
		}()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)

		assert.Equal(t, content, string(data))
	})

	t.Run("distributes content over multiple files", func(t *testing.T) {
		directory, cleanup, err := fsutils.TemporaryDirectory()
		require.NoError(t, err)
		defer cleanup()

		provider := testutils.FilesystemStorageProvider(t, int64(8), directory)

		err = provider.Write("testfile.txt", strings.NewReader("abcdefghijkl"))
		require.NoError(t, err)

		chunkFiles, err := os.ReadDir(filepath.Join(directory, "testfile.txt"))
		require.NoError(t, err)

		require.Equal(t, 2, len(chunkFiles))

		file1Bytes, err := os.ReadFile(filepath.Join(directory, "testfile.txt", chunkFiles[0].Name()))
		assert.NoError(t, err)
		file2Bytes, err := os.ReadFile(filepath.Join(directory, "testfile.txt", chunkFiles[1].Name()))
		assert.NoError(t, err)

		assert.Equal(t, "abcdefgh", string(file1Bytes))
		assert.Equal(t, "ijkl", string(file2Bytes))
	})
}
