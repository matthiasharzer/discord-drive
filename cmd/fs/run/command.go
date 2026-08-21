package run

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/docker/go-units"
	"github.com/matthiasharzer/discord-drive/api"
	"github.com/matthiasharzer/discord-drive/logging"
	"github.com/matthiasharzer/discord-drive/storage/chunk"
	"github.com/matthiasharzer/discord-drive/storage/chunk/filesystem"
	"github.com/matthiasharzer/discord-drive/storage/distributedfile"
	"github.com/matthiasharzer/discord-drive/util/cobrautils"
	"github.com/spf13/cobra"
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
	Use:   "run",
	Short: "Run with a file system storage provider",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if chunkSize.Bytes <= 0 {
			return fmt.Errorf("chunk size must be greater than 0")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		storageProvider := distributedfile.NewStorageProvider(chunkSize.Bytes, func(key string) (chunk.Provider, error) {
			return filesystem.NewChunkProvider(filepath.Join(directory, key)), nil
		})
		defer func() {
			err := storageProvider.Close()
			if err != nil {
				logging.Error("error while closing storage provider", "err", err)
			}
		}()

		mux := api.GetMux(storageProvider)

		addr := fmt.Sprintf("%s:%d", httpHost, httpPort)
		logging.Info("starting discord-drive local mode", "host", httpHost, "port", httpPort)
		err := http.ListenAndServe(
			addr,
			mux,
		)

		return err
	},
}
