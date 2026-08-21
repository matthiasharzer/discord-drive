package run

import (
	"fmt"
	"net/http"
	"os"

	"github.com/matthiasharzer/discord-drive/api"
	"github.com/matthiasharzer/discord-drive/logging"
	"github.com/matthiasharzer/discord-drive/storage/discord"
	"github.com/spf13/cobra"
)

var httpPort int
var httpHost string
var driveChannelID string

func init() {
	Command.Flags().IntVarP(&httpPort, "port", "p", 4000, "The HTTP server port to listen on")
	Command.Flags().StringVarP(&httpHost, "host", "", "", "The HTTP server host (default: all interfaces)")
	Command.Flags().StringVarP(&driveChannelID, "channel", "c", "", "The Discord channel ID to use for storage")
}

var Command = &cobra.Command{
	Use:   "run",
	Short: "Run discord drive",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if driveChannelID == "" {
			return fmt.Errorf("a channel ID must be configured")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		discordToken := os.Getenv("DISCORD_BOT_TOKEN")
		if discordToken == "" {
			return fmt.Errorf("DISCORD_BOT_TOKEN environment variable is not set")
		}

		storageProvider, err := discord.NewStorageProvider(discordToken, driveChannelID)
		if err != nil {
			return fmt.Errorf("failed to create discord storage provider: %w", err)
		}
		defer func() {
			err := storageProvider.Close()
			if err != nil {
				logging.Error("error while closing storage provider", "err", err)
			}
		}()

		mux := api.GetMux(storageProvider)

		addr := fmt.Sprintf("%s:%d", httpHost, httpPort)
		logging.Info("starting discord-drive server", "host", httpHost, "port", httpPort)
		err = http.ListenAndServe(
			addr,
			mux,
		)

		return err
	},
}
