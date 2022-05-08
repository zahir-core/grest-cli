package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	return &cobra.Command{
		Use:     "grest",
		Example: "  grest init",
		Short:   cmdGREST{}.Summary(),
		Long:    cmdGREST{}.Description(),
		Run:     cmdGREST{}.Run,
	}
}

type cmdGREST struct{}

func (cmdGREST) Summary() string {
	return "The command line interface for gREST applications"
}

func (cmdGREST) Description() string {
	return `
The command line interface for gREST applications

Scaffolding GREST project or extension by generating the basic code. The CLI provides the fastest way to get started with a GREST project that adheres to best practices.
`
}

func (cmdGREST) Example() string {
	return "  grest init"
}

func (cmdGREST) Run(c *cobra.Command, args []string) {
	isPrintVersion, _ := c.Flags().GetBool("version")
	if isPrintVersion {
		PrintVersion()
	} else {
		c.Help()
	}
}
