package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type PromptType int

const (
	TextPrompt     PromptType = 0
	PasswordPrompt PromptType = 1
	SelectPrompt   PromptType = 2
)

type promptItem struct {
	ID            string
	Label         string
	Value         string
	SelectOptions []string
	promptType    PromptType
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings for the application",
	Long:  `Configure settings for the application`,
	Run: func(cmd *cobra.Command, args []string) {
		items := []*promptItem{
			{
				ID:         "APIKey",
				Label:      "API Key",
				promptType: PasswordPrompt,
			},
			{
				ID:            "Theme",
				Label:         "Theme",
				promptType:    SelectPrompt,
				SelectOptions: []string{"Dark", "Light"},
			},
			{
				ID:            "Language",
				Label:         "Preferred Language",
				promptType:    SelectPrompt,
				SelectOptions: []string{"English", "Spanish", "French", "German", "Chinese", "Japanese"},
			},
		}

		promptNested("Configuration Items", 0, items)

		fmt.Println("")
		for _, v := range items {
			fmt.Printf("Saving configuration (%s) with value (%s)...\n", v.ID, v.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func promptNested(promptLabel string, startingIndex int, items []*promptItem) bool {

	doneID := "Done"

	if len(items) > 0 && items[0].ID != doneID {
		items = append([]*promptItem{{ID: doneID, Label: "Done"}}, items...)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Label | cyan }}",
		Inactive: "{{ .Label | cyan }}",
		Selected: "\U0001F336 {{ .Label | red  | cyan }}",
	}

	prompt := promptui.Select{
		Label:        promptLabel,
		Items:        items,
		Templates:    templates,
		Size:         3,
		HideSelected: true,
		CursorPos:    startingIndex,
	}

	idx, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Error occured when running prompt: %v\n", err)
		return false
	}

	selectedItem := items[idx]

	if selectedItem.ID == doneID {
		return true
	}

	var promptResponse string

	if selectedItem.promptType == TextPrompt || selectedItem.promptType == PasswordPrompt {
		promptResponse, err = promptInput(*selectedItem)

		if err != nil {
			fmt.Printf("Error occured when running prompt: %v\n", err)
			return false
		}

		items[idx].Value = promptResponse

	}

	if selectedItem.promptType == SelectPrompt {
		promptResponse, err = promptSelect(*selectedItem)

		if err != nil {
			fmt.Printf("Error occured when running prompt: %v\n", err)
			return false
		}
		items[idx].Value = promptResponse
	}

	if err != nil {
		fmt.Printf("Error occured when running prompt: %v\n", err)
		return false
	}

	return promptNested(promptLabel, idx, items)
}

func promptInput(item promptItem) (string, error) {
	prompt := promptui.Prompt{
		Label:       item.Label,
		HideEntered: true,
	}

	if item.promptType == PasswordPrompt {
		prompt.Mask = '*'
	}

	res, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return res, nil
}

func promptSelect(item promptItem) (string, error) {
	prompt := promptui.Select{
		Label:        item.Label,
		Items:        item.SelectOptions,
		HideSelected: true,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}
