package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Genera el script de autocompletado para tu shell",
	Long: `Genera autocompletado para cliweather.

Bash:
  source <(cliweather completion bash)
  # persistente:
  cliweather completion bash > /etc/bash_completion.d/cliweather  (root)
  # o: cliweather completion bash > ~/.local/share/bash-completion/cliweather

Zsh:
  echo 'autoload -U compinit; compinit' >> ~/.zshrc
  cliweather completion zsh > "${fpath[1]}/_cliweather"   # requiere $fpath writable
  # o: mkdir -p ~/.zsh/completions && cliweather completion zsh > ~/.zsh/completions/_cliweather
  #    y añade a ~/.zshrc: fpath=(~/.zsh/completions $fpath)

Fish:
  cliweather completion fish > ~/.config/fish/completions/cliweather.fish

PowerShell:
  cliweather completion powershell | Out-String | Invoke-Expression
  # persistente:
  cliweather completion powershell > $PROFILE
`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{
		"bash", "zsh", "fish", "powershell",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]
		switch shell {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			// Cobra recomienda head para zsh
			fmt.Fprintln(os.Stdout, "#compdef cliweather")
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s", shell)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
