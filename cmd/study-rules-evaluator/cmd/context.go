package cmd

import (
  "fmt"
  "log"
  "github.com/spf13/cobra"
  "github.com/grippenet/study-rules-evaluator/evaluator/rules"

)

func init() {
  rootCmd.AddCommand(contextCmd)
}

var contextCmd = &cobra.Command{
  Use:   "context",
  Short: "Find Context of study rules",
  Long:  `Show rules context, event type and survey subject of the rules`,
  Run: func(cmd *cobra.Command, args []string) {
    studyRules, err := loadStudyRules()
    if(err != nil) {
      log.Fatalf("Unable to read rules : %s", err)
    }

	for i, rule := range studyRules {
		ctx := rules.FindRuleContext(rule)
		fmt.Printf("Rule %d %s\n", i, ctx.EventType)
		for _, sc := range ctx.Surveys  {
			fmt.Printf(" - %s\n", sc)
		}
	}
  },
}