package server

import (
	"github.com/matthiasharzer/discord-drive/cmd/server/run"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(run.Command)
}

var Command = &cobra.Command{
	Use: "server",
}
