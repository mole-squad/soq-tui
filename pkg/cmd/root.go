package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mole-squad/soq-tui/pkg/app"
	"github.com/spf13/cobra"
)

const (
	debugFlagKey     = "debug"
	configDirFlagKey = "config-dir"
)

var debugEnabled bool

var rootCmd = &cobra.Command{
	Use: "qt",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool(debugFlagKey)
		configDir, _ := cmd.Flags().GetString(configDirFlagKey)

		m := app.NewAppModel(
			app.WithDebugMode(debug),
			app.WithConfigDir(configDir),
		)

		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home directory: %v\n", err)
		os.Exit(1)
	}

	defaultConfigDir := filepath.Join(homeDir, ".soq")

	rootCmd.PersistentFlags().BoolP(debugFlagKey, "d", false, "enable debug mode")
	rootCmd.PersistentFlags().StringP(configDirFlagKey, "c", defaultConfigDir, "config directory")
}
