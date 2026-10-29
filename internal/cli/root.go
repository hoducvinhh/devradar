package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/devradar/devradar/internal/app"
	"github.com/spf13/cobra"
)

var (
	flagToken   string
	flagFormat  string
	flagVerbose bool

	application *app.App
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "devradar",
	Short: "DevRadar – Find your next Open Source contribution",
	Long: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6366f1")).Render(`
  ____             ____           _             
 |  _ \  _____   _|  _ \ __ _  __| | __ _ _ __  
 | | | |/ _ \ \ / / |_) / _` + "`" + ` |/ _` + "`" + ` |/ _` + "`" + ` | '__| 
 | |_| |  __/\ V /|  _ < (_| | (_| | (_| | |    
 |____/ \___| \_/ |_| \_\__,_|\__,_|\__,_|_|    
`) + `
DevRadar connects you with open-source issues matching your skills,
evaluates issue difficulty, checks repository health, and explains why each issue is recommended.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		application, err = app.NewApp(flagToken)
		if err != nil {
			return fmt.Errorf("initializing application: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "GitHub Personal Access Token")
	RootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "table", "Output format (table, json, csv)")
	RootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output")

	// Register subcommands
	RootCmd.AddCommand(findCmd)
	RootCmd.AddCommand(viewCmd)
	RootCmd.AddCommand(profileCmd)
	RootCmd.AddCommand(saveCmd)
	RootCmd.AddCommand(savedCmd)
	RootCmd.AddCommand(unsaveCmd)
	RootCmd.AddCommand(authCmd)
	RootCmd.AddCommand(repoCmd)
}

func printRateLimitWarning() {
	if application != nil && application.Client != nil {
		if msg := application.Client.RateLimit().WarningMessage(); msg != "" {
			fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Render(msg))
		}
	}
}
