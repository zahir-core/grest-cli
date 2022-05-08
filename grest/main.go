package main

import (
	"github.com/spf13/cobra"

	"grest.dev/cmd"
)

var boolVar = false

func main() {
	cobra.EnableCommandSorting = false

	cli := cmd.New()
	cli.AddCommand(cmd.CmdInit())
	cli.AddCommand(cmd.CmdAdd())
	cli.AddCommand(cmd.CmdFmt())
	cli.AddCommand(cmd.CmdVersion())

	cli.CompletionOptions.DisableDefaultCmd = true
	cli.PersistentFlags().BoolVarP(&boolVar, "version", "v", false, "Print the grest version")
	cli.Flags().SortFlags = false
	cli.Execute()
}
