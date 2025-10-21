package cliweather

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Shell script",
	Long: `It generates an auto-completed script for your shell.

Bash:
	source <(cliweather completion bash)
	# persistente:
	cliweather completion bash > /etc/bash_completion.d/cliweather  (root)
	# o: cliweather completion bash > ~/.local/share/bash-completion/cliweather

Zsh:
	echo 'autoload -U compinit; compinit' >> ~/.zshrc
	cliweather completion zsh > "${fpath[1]}/_cliweather"   # requires $fpath writable
	# o: mkdir -p ~/.zsh/completions && cliweather completion zsh > ~/.zsh/completions/_cliweather
	#    and adds ~/.zshrc: fpath=(~/.zsh/completions $fpath)

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
			fmt.Fprintln(os.Stdout, "#compdef cliweather")
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("Unsupported shell: %s", shell)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
