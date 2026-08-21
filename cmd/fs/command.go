package fs

import (
	"github.com/matthiasharzer/discord-drive/cmd/fs/run"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(run.Command)
}

var Command = &cobra.Command{
	Use:   "fs",
	Short: "Run with a local filesystem storage provider",
}
