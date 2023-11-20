package cmd

import (
	"os"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

// https://github.com/manifoldco/promptui/issues/49#issuecomment-428801411
type stderr struct{}

func (s *stderr) Write(b []byte) (int, error) {
	if len(b) == 1 && b[0] == 7 {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (s *stderr) Close() error {
	return os.Stderr.Close()
}

var rootCmd = &cobra.Command{
	Use:   "nested-prompt",
	Short: "A CLI application that demonstrates how to use nested prompts in Cobra using promptui",
	Long:  `A CLI application that demonstrates how to use nested prompts in Cobra using promptui.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	readline.Stdout = &stderr{}
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
