package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	noColor bool
	noEmoji bool
)

var rootCmd = &cobra.Command{
	Use:           "cliweather",
	Short:         "Easy to use weather CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable ANSI colors")
	rootCmd.PersistentFlags().BoolVar(&noEmoji, "no-emoji", false, "Disable emojis")
}

func envNoColor() bool {
	// NO_COLOR standard: https://no-color.org/
	return os.Getenv("NO_COLOR") != ""
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
