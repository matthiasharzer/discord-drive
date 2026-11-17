package storagetest

import (
	"io"
	"strings"
	"testing"

	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunProviderSuite(t *testing.T, newProvider func(t *testing.T) storage.Provider) {
	t.Run("reads a clean divided file", func(t *testing.T) {
		provider := newProvider(t)

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

	t.Run("reads large saved file", func(t *testing.T) {
		provider := newProvider(t)

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

}
