package tools

import "github.com/spf13/cobra"

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
	cmd := &cobra.Command{
		Use:   "new <projectName>",
		Short: "Create a new HandyGo project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewProject(ProjectNewOptions{
				ProjectName: args[0],
				ModulePath:  modulePath,
				Force:       force,
			})
		},
	}
	cmd.Flags().StringVarP(&modulePath, "module", "m", "", "new project module path")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite scaffold files when the target directory exists")
	return cmd
}
