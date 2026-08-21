package api

import (
	"net/http"

	"github.com/matthiasharzer/discord-drive/api/file"
	"github.com/matthiasharzer/discord-drive/api/upload"
	"github.com/matthiasharzer/discord-drive/storage"
)

func GetMux(storageProvider storage.Provider) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/v1/upload/{identifier}", upload.Handler(storageProvider))
	mux.HandleFunc("GET /api/v1/file/{identifier}", file.Handler(storageProvider))

	return mux
}
