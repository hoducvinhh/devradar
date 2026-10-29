package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/devradar/devradar/internal/output"
	"github.com/devradar/devradar/internal/profile"
	"github.com/spf13/cobra"
)

var (
	profSetLevel     string
	profSetLanguages string
	profSetInterests string
	profSetMinStars  int
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage local developer profile",
	Long:  `Manage your skills, experience level, and interests to personalize issue rankings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProfileShow(cmd, args)
	},
}

var profileShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current profile settings",
	RunE:  runProfileShow,
}

func runProfileShow(cmd *cobra.Command, args []string) error {
	p, err := application.Profile.Get()
	if err != nil {
		return fmt.Errorf("loading profile: %w", err)
	}

	if flagFormat == "json" {
		return output.RenderJSON(os.Stdout, p)
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6366f1"))
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#94a3b8"))
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f8fafc"))

	fmt.Println(headerStyle.Render("\n👤 DevRadar Developer Profile (~/.devradar/profile.json)"))
	fmt.Printf("%s: %s\n", labelStyle.Render("Target Level"), valStyle.Render(string(p.Level)))

	langs := "None"
	if len(p.Languages) > 0 {
		langs = strings.Join(p.Languages, ", ")
	}
	fmt.Printf("%s: %s\n", labelStyle.Render("Languages"), valStyle.Render(langs))

	interests := "None"
	if len(p.Interests) > 0 {
		interests = strings.Join(p.Interests, ", ")
	}
	fmt.Printf("%s: %s\n", labelStyle.Render("Interests"), valStyle.Render(interests))

	if p.MinStars > 0 {
		fmt.Printf("%s: %s\n", labelStyle.Render("Min Stars"), valStyle.Render(fmt.Sprintf("%d", p.MinStars)))
	}

	fmt.Printf("%s: %s\n\n", labelStyle.Render("Last Updated"), valStyle.Render(p.UpdatedAt.Format("2006-01-02 15:04:05")))
	fmt.Printf("Tip: Update profile using 'devradar profile set --level beginner --languages go,python'\n\n")

	return nil
}

var profileSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Update profile settings",
	Example: `  devradar profile set --level beginner --languages go,python --interests backend,cli`,
	RunE: func(cmd *cobra.Command, args []string) error {
		langs := profile.ParseCommaSeparated(profSetLanguages)
		interests := profile.ParseCommaSeparated(profSetInterests)

		updated, err := application.Profile.Update(profSetLevel, langs, interests, profSetMinStars)
		if err != nil {
			return fmt.Errorf("saving profile: %w", err)
		}

		if flagFormat == "json" {
			return output.RenderJSON(os.Stdout, updated)
		}

		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render("✔ Profile updated successfully!"))
		return runProfileShow(cmd, args)
	},
}

var profileResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset profile to default anonymous settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := application.Profile.Reset(); err != nil {
			return fmt.Errorf("resetting profile: %w", err)
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render("✔ Profile reset to default settings."))
		return nil
	},
}

func init() {
	profileSetCmd.Flags().StringVar(&profSetLevel, "level", "", "Target difficulty level (beginner, intermediate, advanced)")
	profileSetCmd.Flags().StringVar(&profSetLanguages, "languages", "", "Comma-separated preferred languages (e.g. go,python,rust)")
	profileSetCmd.Flags().StringVar(&profSetInterests, "interests", "", "Comma-separated topics (e.g. backend,cli,kubernetes)")
	profileSetCmd.Flags().IntVar(&profSetMinStars, "min-stars", -1, "Preferred minimum repository stars")

	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileSetCmd)
	profileCmd.AddCommand(profileResetCmd)
}
