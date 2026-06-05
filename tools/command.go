package tools

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "handygo",
		Short:         "HandyGo development tools",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.AddCommand(newNewCommand())
	rootCmd.AddCommand(newGenCommand())
	return rootCmd
}

func newGenCommand() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate project code",
	}
	genCmd.AddCommand(newGenModelCommand())
	return genCmd
}

func newGenModelCommand() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:   "model",
		Short: "Generate Gorm models from database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			return GenerateModels(ModelGenerateOptions{ConfigPath: configPath})
		},
	}
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "generator config file")
	return cmd
}

func newNewCommand() *cobra.Command {
	var modulePath string
	var force bool
	var frontend bool
	var noFrontend bool
	var frontendStyleSkill string
	var skipTidy bool
	cmd := &cobra.Command{
		Use:   "new <projectName>",
		Short: "Create a new HandyGo project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if frontend && noFrontend {
				return fmt.Errorf("--frontend and --no-frontend cannot be used together")
			}
			var frontendOpt *bool
			if frontend || noFrontend {
				value := frontend && !noFrontend
				frontendOpt = &value
			}
			return NewProject(ProjectNewOptions{
				ProjectName:        args[0],
				ModulePath:         modulePath,
				Force:              force,
				Interactive:        frontendOpt == nil && isTerminalInput(),
				Frontend:           frontendOpt,
				FrontendStyleSkill: frontendStyleSkill,
				SkipTidy:           skipTidy,
			})
		},
	}
	cmd.Flags().StringVarP(&modulePath, "module", "m", "", "new project module path")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite scaffold files when the target directory exists")
	cmd.Flags().BoolVar(&frontend, "frontend", false, "enable frontend agent workflow")
	cmd.Flags().BoolVar(&noFrontend, "no-frontend", false, "disable frontend agent workflow")
	cmd.Flags().StringVar(&frontendStyleSkill, "frontend-style-skill", "", "default frontend style skill name")
	cmd.Flags().BoolVar(&skipTidy, "skip-tidy", false, "skip running go mod tidy after project creation")
	return cmd
}

func isTerminalInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
