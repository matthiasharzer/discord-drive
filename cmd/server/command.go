package server

import (
	"github.com/spf13/cobra"

	"github.com/matthiasharzer/discord-drive/cmd/server/run"
)

func init() {
	Command.AddCommand(run.Command)
}

var Command = &cobra.Command{
	Use: "server",
}
