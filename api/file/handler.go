package file

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/matthiasharzer/discord-drive/storage"
)

func Handler(storageProvider storage.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		identifier := r.PathValue("identifier")
		if identifier == "" {
			http.Error(w, "missing file identifier", http.StatusBadRequest)
			return
		}

		if identifier == "." || identifier == ".." || identifier != filepath.Base(identifier) || strings.ContainsAny(identifier, `/\\`) {
			http.Error(w, "invalid file identifier", http.StatusBadRequest)
			return
		}

		exists, err := storageProvider.Has(identifier)
		if err != nil {
			http.Error(w, "failed to retrieve file", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		reader, err := storageProvider.Read(identifier)
		if err != nil {
			http.Error(w, "failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = reader.Close()
		}()

		_, err = io.Copy(w, reader)
		if err != nil {
			http.Error(w, "failed to retrieve file", http.StatusInternalServerError)
			return
		}
	}
}
