package upload_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matthiasharzer/discord-drive/api/upload"
	"github.com/matthiasharzer/discord-drive/util/testutils"
)

func TestHandler(t *testing.T) {
	t.Run("uploads a file", func(t *testing.T) {
		provider, cleanup := testutils.TempDirFilesystemStorageProvider(t, int64(8))
		defer cleanup()

		request, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("abcdefghijkl"))
		require.NoError(t, err)
		request.SetPathValue("identifier", "testfile.txt")

		recorder := httptest.NewRecorder()
		handler := upload.Handler(provider)
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "", recorder.Body.String())

		exists, err := provider.Has("testfile.txt")
		require.NoError(t, err)
		assert.True(t, exists)

		reader, err := provider.Read("testfile.txt")
		require.NoError(t, err)
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "abcdefghijkl", string(data))
	})

	t.Run("fails if HTTP method other than POST is used", func(t *testing.T) {
		disallowedMethods := []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		}
		for _, method := range disallowedMethods {
			t.Run("fails with method "+method, func(t *testing.T) {
				request, err := http.NewRequest(method, "/", nil)
				request.SetPathValue("identifier", "testfile.txt")
				require.NoError(t, err)

				recorder := httptest.NewRecorder()

				handler := upload.Handler(nil)
				handler.ServeHTTP(recorder, request)

				assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
				assert.Equal(t, "method not allowed\n", recorder.Body.String())
			})
		}
	})
}
