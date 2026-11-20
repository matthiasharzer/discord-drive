package api

import (
	"net/http"

	"github.com/matthiasharzer/discord-drive/api/file"
	"github.com/matthiasharzer/discord-drive/api/upload"
	"github.com/matthiasharzer/discord-drive/storage"
	"github.com/matthiasharzer/discord-drive/ui"

	"github.com/matthiasharzer/service-kit/util/httputil"
)

func GetMux(storageProvider storage.Provider) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/v1/upload/{identifier}", upload.Handler(storageProvider))
	mux.HandleFunc("GET /api/v1/file/{identifier}", file.Handler(storageProvider))
	mux.Handle("GET /api/", http.NotFoundHandler())

	mux.Handle("GET /",
		httputil.UseMiddleware(
			[]httputil.Middleware{httputil.GZIPMiddleware()},
			httputil.HandleStaticSite(ui.Content),
		),
	)

	return mux
}
