package cmd

import (
	"strings"
	"fmt"
	"os"

	"github.com/mgutz/ansi"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type questionnaire struct {
	template *survey.Select
	projDir *survey.Input
}

var (
	// Used for cli arguments.
	template string
	projDir  string

	templates = map[string]string{
		"gorsk": "https://github.com/ribice/gorsk/archive/master.zip",
	}

	validTemplates = keys(templates)

	rootCmd = &cobra.Command{
		Use:     "goinit [template] [project-directory]",
		Short:   "A generator for my go templates",
		Args:    cobra.MaximumNArgs(2),
		Example: "  goinit go-kit-microservice my-microservice",
		Run:     run,
	}

	questions = questionnaire{
		template: &survey.Select{
			Message: "Which template do you want to use?",
			Options: validTemplates,
		},
		projDir: &survey.Input{
			Message: "Specify the project directory:",
			Default: "my-app",
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		survey.AskOne(questions.template, &template, survey.WithValidator(survey.Required))
	} else {
		template = args[0]
	}
	if len(args) < 2 {
		survey.AskOne(questions.projDir, &projDir, survey.WithValidator(survey.Required))
	} else {
		projDir = args[1]
	}

	checkValidArgs(template, projDir)
	os.MkdirAll(projDir, 0777)

	println("template", template)
	println("projDir", projDir)
}

func checkValidArgs(template string, projDir string) {
	if template == "" || projDir == "" {
		os.Exit(1)
	}

	coloredTmpls := make([]string, 0, len(templates))
	for key := range templates {
		coloredTmpls = append(coloredTmpls, ansi.Color(key, "magenta+bh"))
		if key == template {
			return
		}
	}

	fmt.Println("Error: invalid template")
	fmt.Printf("Available templates: %s\n", strings.Join(coloredTmpls, ", "))
}

func keys(data map[string]string) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
