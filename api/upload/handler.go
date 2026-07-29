package upload

import (
	"net/http"

	"github.com/matthiasharzer/discord-drive/storage"

	"github.com/docker/go-units"
)

const fileSizeLimit = 10 * units.GiB

func Handler(storageProvider storage.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		identifier := r.PathValue("identifier")
		if identifier == "" {
			http.Error(w, "missing file identifier", http.StatusBadRequest)
			return
		}
		if r.Body == nil {
			http.Error(w, "missing request body", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		limitedReader := http.MaxBytesReader(w, r.Body, int64(fileSizeLimit))

		err := storageProvider.Write(identifier, limitedReader)
		if err != nil {
			http.Error(w, "failed to store file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
