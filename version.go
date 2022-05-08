package cmd

import (
	"fmt"
	"strings"

	"grest.dev/grest/log"
)

var Version = "v0.0.0"

func PrintVersion() {
	msg := strings.Builder{}
	msg.WriteString(log.Fmt(`        ________________________________________`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`       /        `, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString(log.Fmt(`____`, log.HiCyan, log.Bold, log.BlinkRapid))
	msg.WriteString(log.Fmt(`___  `, log.Red, log.Bold))
	msg.WriteString(log.Fmt(`____`, log.Yellow, log.Bold))
	msg.WriteString(log.Fmt(`____`, log.Green, log.Bold))
	msg.WriteString(log.Fmt(`_____ `, log.Blue, log.Bold))
	msg.WriteString(log.Fmt(`       /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`      /    `, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString(log.Fmt(`--- / __/`, log.HiCyan, log.Bold, log.BlinkRapid))
	msg.WriteString(log.Fmt(` _ \`, log.Red, log.Bold))
	msg.WriteString(log.Fmt(`/ __/`, log.Yellow, log.Bold))
	msg.WriteString(log.Fmt(` __/`, log.Green, log.Bold))
	msg.WriteString(log.Fmt(`_  _/ `, log.Blue, log.Bold))
	msg.WriteString(log.Fmt(`      /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`     /   `, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString(log.Fmt(`---- / / /`, log.HiCyan, log.Bold, log.BlinkRapid))
	msg.WriteString(log.Fmt(` / _/`, log.Red, log.Bold))
	msg.WriteString(log.Fmt(` _/`, log.Yellow, log.Bold))
	msg.WriteString(log.Fmt(`_\ \`, log.Green, log.Bold))
	msg.WriteString(log.Fmt(`  / /  `, log.Blue, log.Bold))
	msg.WriteString(log.Fmt(`      /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`    /     `, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString(log.Fmt(`-- /___/`, log.HiCyan, log.Bold, log.BlinkRapid))
	msg.WriteString(log.Fmt(`_/\ \`, log.Red, log.Bold))
	msg.WriteString(log.Fmt(`___/`, log.Yellow, log.Bold))
	msg.WriteString(log.Fmt(`___/`, log.Green, log.Bold))
	msg.WriteString(log.Fmt(` /_/ `, log.Blue, log.Bold))
	msg.WriteString(log.Fmt(`CLI`, log.HiCyan, log.Bold))
	msg.WriteString(log.Fmt(`    /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`   /                                      /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`  /               `, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString(log.Fmt(" ", log.BgRed))
	msg.WriteString(log.Fmt(Version, log.BgRed, log.Bold))
	msg.WriteString(log.Fmt(" ", log.BgRed))
	msg.WriteString(log.Fmt(`               /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(` /             `, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString(log.Fmt("https://grest.dev", log.Blue))
	msg.WriteString(log.Fmt(`        /`, log.HiMagenta, log.Bold, log.Italic))
	msg.WriteString("\n")

	msg.WriteString(log.Fmt(`/______________________________________/`, log.HiMagenta, log.Bold, log.Italic))

	fmt.Fprintln(log.Stdout, msg.String())
}
