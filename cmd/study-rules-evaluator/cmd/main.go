package cmd

import (
	"fmt"
	"os"
	"github.com/grippenet/study-rules-evaluator/evaluator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "study-rules-evaluator",
	Short: "Influenzanet Study-rules-evaluator",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadStudyRules() (evaluator.StudyRules, error) {
	file := os.Getenv("STUDYRULES_FILE")
	if(rulesFile != "") {
		file = rulesFile
	}
	if(file == "") {
		fmt.Println("Study rules file is not specified, use --rules or `STUDYRULES_FILE` environment variable")
		os.Exit(2)
	}
	return evaluator.ReadRulesFromJSON(file)
}

var rulesFile string


func init() {
	scenarioCmd.Flags().StringVar(&rulesFile, "rules", "", "Rules file")
}
