package file_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matthiasharzer/discord-drive/api/file"
	"github.com/matthiasharzer/discord-drive/util/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("fails if HTTP method other than GET is used", func(t *testing.T) {
		disallowedMethods := []string{
			http.MethodPost,
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

				handler := file.Handler(nil)
				handler.ServeHTTP(recorder, request)

				assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
				assert.Equal(t, "method not allowed\n", recorder.Body.String())
			})
		}
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
