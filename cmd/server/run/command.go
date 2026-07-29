package run

import (
	"fmt"
	"net/http"
	"path"

	"github.com/docker/go-units"
	"github.com/matthiasharzer/discord-drive/util/cobrautils"
	"github.com/spf13/cobra"

	"github.com/matthiasharzer/discord-drive/api/file"
	"github.com/matthiasharzer/discord-drive/api/upload"
	"github.com/matthiasharzer/discord-drive/logging"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/filesystem"
	"github.com/matthiasharzer/discord-drive/storage/distributedfiles"
)

var httpPort int
var httpHost string
var chunkSize = cobrautils.FlagFileSize{Bytes: 100 * units.MB}

func init() {
	Command.Flags().IntVarP(&httpPort, "port", "p", 4000, "HTTP server port")
	Command.Flags().StringVarP(&httpHost, "host", "", "", "HTTP server host (default: all interfaces)")
	Command.Flags().VarP(&chunkSize, "chunk-size", "", "chunk size")
}

var Command = &cobra.Command{
	Use: "run",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if chunkSize.Bytes <= 0 {
			return fmt.Errorf("chunk size must be greater than 0")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		storageProvider := distributedfiles.NewProvider(int64(2*units.KiB), func(key string) chunk.Provider {
			return filesystem.NewProvider(path.Join(".test", "data", key))
		})

		mux.HandleFunc("POST /api/v1/upload/{identifier}", upload.Handler(storageProvider))
		mux.HandleFunc("GET /api/v1/file/{identifier}", file.Handler(storageProvider))

		addr := fmt.Sprintf("%s:%d", httpHost, httpPort)
		logging.Info("starting sync-watch-server", "host", httpHost, "port", httpPort)
		err := http.ListenAndServe(
			addr,
			mux,
		)

		return err
	},
}
