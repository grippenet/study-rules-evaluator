package cmd

import (
  "fmt"
  "log"
  "github.com/spf13/cobra"
   "github.com/grippenet/study-rules-evaluator/evaluator/scenario"
  "github.com/grippenet/study-rules-evaluator/evaluator/engine"
)

func init() {
  rootCmd.AddCommand(scenarioCmd)
  scenarioCmd.Flags().StringVar(&scenarioFile, "file", "", " Scenario file")
}

var scenarioFile string

var scenarioCmd = &cobra.Command{
  Use:   "scenario",
  Short: "Run and evaluate a submt scenario",
  Long:  `Evaluate a submit scenario`,
  Run: func(cmd *cobra.Command, args []string) {
    studyRules, err := loadStudyRules()
    if(err != nil) {
      log.Fatalf("Unable to read rules : %s", err)
    }

    scenarios, err := scenario.Load(scenarioFile)
    if(err != nil) {
      fmt.Println("Error loading scenarios in %s :", scenarioFile,  err)
      return
    }
    
    for idx, scenario := range scenarios {
      evaluator := engine.NewRuleEvaluator(nil, studyRules)
      //evaluator.Verbose = true
      result := scenario.Run(evaluator)
      fmt.Printf("Scenario %d %s\n", idx, scenario.Label)
      scenario.PrintResult(result)  
    }
    
  },
}