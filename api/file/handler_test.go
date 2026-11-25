package file_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matthiasharzer/discord-drive/api/file"
	"github.com/matthiasharzer/discord-drive/util/testutils"
)

func TestHandler(t *testing.T) {
	t.Run("returns the exact same content as stored", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(8))
		defer cleanup()

		err := provider.Write("testfile.txt", strings.NewReader("abcdefghijkl"))
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		request.SetPathValue("identifier", "testfile.txt")
		require.NoError(t, err)

		recorder := httptest.NewRecorder()

		handler := file.Handler(provider)
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "abcdefghijkl", recorder.Body.String())
	})

	t.Run("fails if unknown identifier is used", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(8))
		defer cleanup()

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		request.SetPathValue("identifier", "testfile.txt")
		require.NoError(t, err)

		recorder := httptest.NewRecorder()

		handler := file.Handler(provider)
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, "file not found\n", recorder.Body.String())
	})
}
