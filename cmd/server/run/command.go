package run

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/docker/go-units"
	"github.com/spf13/cobra"

	"github.com/matthiasharzer/discord-drive/api/file"
	"github.com/matthiasharzer/discord-drive/api/upload"
	"github.com/matthiasharzer/discord-drive/logging"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/filesystem"
	"github.com/matthiasharzer/discord-drive/storage/distributedfile"
	"github.com/matthiasharzer/discord-drive/util/cobrautils"
)

var httpPort int
var httpHost string
var chunkSize = cobrautils.FlagFileSize{Bytes: 100 * units.MB}
var directory string

func init() {
	Command.Flags().IntVarP(&httpPort, "port", "p", 4000, "The HTTP server port to listen on")
	Command.Flags().StringVarP(&httpHost, "host", "", "", "The HTTP server host (default: all interfaces)")
	Command.Flags().VarP(&chunkSize, "chunk-size", "", "The maximum size of a single file chunk (e.g. 100MB)")
	Command.Flags().StringVarP(&directory, "directory", "d", ".", "The directory to store files in")
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

		storageProvider := distributedfile.NewProvider(chunkSize.Bytes, func(key string) chunk.Provider {
			return filesystem.NewProvider(filepath.Join(directory, key))
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
