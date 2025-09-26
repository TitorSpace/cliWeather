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
	Short:         "CLI del tiempo sencilla y práctica",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags persistentes disponibles para todos los subcomandos
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Desactivar colores ANSI en la salida")
	rootCmd.PersistentFlags().BoolVar(&noEmoji, "no-emoji", false, "Desactivar emojis en la salida")
}

// ===== Helpers de entorno para color/emoji =====

func envNoColor() bool {
	// Estándar NO_COLOR: https://no-color.org/
	return os.Getenv("NO_COLOR") != ""
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	// Si es un dispositivo de carácter, asumimos TTY
	return (fi.Mode() & os.ModeCharDevice) != 0
}
