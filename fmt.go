package cmd

import (
	"github.com/spf13/cobra"
	"grest.dev/grest"
)

var (
	// main operation modes
	list        = false
	write       = true
	rewriteRule = ""
	simplifyAST = false
	doDiff      = false
	allErrors   = false

	// debugging
	cpuprofile = ""
)

type cmdFmt struct{}

func CmdFmt() *cobra.Command {
	cli := &cobra.Command{
		Use:   "fmt [flags] [path ...]",
		Short: cmdFmt{}.Summary(),
		Long:  cmdFmt{}.Description(),
		Run:   cmdFmt{}.Run,
	}
	// cli.Flags().BoolVarP(&list, "list", "l", false, "list files whose formatting differs from gofmt's")
	cli.Flags().BoolVarP(&write, "write", "w", true, "write result to (source) file instead of stdout")
	// cli.Flags().StringVarP(&rewriteRule, "rewrite-rule", "r", "", "rewrite rule (e.g., 'a[b:len(a)] -> a[b:]')")
	// cli.Flags().BoolVarP(&simplifyAST, "simplify", "s", false, "simplify code")
	// cli.Flags().BoolVarP(&doDiff, "display-diffs", "d", false, "display diffs instead of rewriting files")
	// cli.Flags().BoolVarP(&allErrors, "report-errors", "e", false, "report all errors (not just the first 10 on different lines)")

	// cli.Flags().StringVarP(&cpuprofile, "cpuprofile", "", "", "write cpu profile to this file")

	return cli
}

func (cmdFmt) Summary() string {
	return "Format the struct tag"
}

func (cmdFmt) Description() string {
	return `
Format / align the multiple struct tags. It align the tags to become much more readable and save allot of time from having to do it manually.

Ensure you run this within the root directory of your app.
`
}

func (cmdFmt) Run(c *cobra.Command, args []string) {
	grest.FormatFile(args...)
}
