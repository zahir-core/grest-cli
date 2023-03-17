package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"grest.dev/grest"
)

const Version = "v0.0.8"

type cmdVersion struct{}

func CmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: cmdVersion{}.Summary(),
		Long:  cmdVersion{}.Description(),
		Run:   cmdVersion{}.Run,
	}
}

func (cmdVersion) Summary() string {
	return "Print the grest version"
}

func (cmdVersion) Description() string {
	return `
Print the grest version.
`
}

func (cmdVersion) Run(c *cobra.Command, args []string) {
	PrintVersion()
}

func PrintVersion() {
	fmt.Fprintln(grest.FmtStdout, GetVersion())
}

func GetVersion() string {
	msg := strings.Builder{}
	msg.WriteString(grest.Fmt(`        ________________________________________`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`       /        `, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString(grest.Fmt(`____`, grest.FmtHiCyan, grest.FmtBold, grest.FmtBlinkRapid))
	msg.WriteString(grest.Fmt(`___  `, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`____`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`____`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`_____ `, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`       /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`      /    `, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString(grest.Fmt(`--- / __/`, grest.FmtHiCyan, grest.FmtBold, grest.FmtBlinkRapid))
	msg.WriteString(grest.Fmt(` _ \`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`/ __/`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(` __/`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`_  _/ `, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`      /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`     /   `, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString(grest.Fmt(`---- / / /`, grest.FmtHiCyan, grest.FmtBold, grest.FmtBlinkRapid))
	msg.WriteString(grest.Fmt(` / _/`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(` _/`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`_\ \`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`  / /  `, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`      /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`    /     `, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString(grest.Fmt(`-- /___/`, grest.FmtHiCyan, grest.FmtBold, grest.FmtBlinkRapid))
	msg.WriteString(grest.Fmt(`_/\ \`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`___/`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`___/`, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(` /_/ `, grest.FmtBlue, grest.FmtBold))
	msg.WriteString(grest.Fmt(`CLI`, grest.FmtHiCyan, grest.FmtBold))
	msg.WriteString(grest.Fmt(`    /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`   /                                      /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`  /               `, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString(grest.Fmt(" ", grest.FmtBgRed))
	msg.WriteString(grest.Fmt(Version, grest.FmtBgRed, grest.FmtBold))
	msg.WriteString(grest.Fmt(" ", grest.FmtBgRed))
	msg.WriteString(grest.Fmt(`               /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(` /             `, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString(grest.Fmt("https://grest.dev", grest.FmtBlue))
	msg.WriteString(grest.Fmt(`        /`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	msg.WriteString("\n")

	msg.WriteString(grest.Fmt(`/______________________________________/`, grest.FmtHiMagenta, grest.FmtBold, grest.FmtItalic))
	return msg.String()
}
