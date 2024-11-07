package cmd

import (
  "os"
  "fmt"
  "errors"
  "github.com/spf13/cobra"
  "github.com/grippenet/study-rules-evaluator/evaluator"
  "github.com/grippenet/study-rules-evaluator/evaluator/scenario"
)

func init() {
  // Default get from env
  rootCmd.AddCommand(scenarioCmd)
  scenarioCmd.Flags().StringVar(&scenarioFile, "file", "", "Scenario file")
  scenarioCmd.Flags().StringVar(&externalServicesFile, "externals", os.Getenv("EXTERNAL_SERVICES_FILE"), "External services definition file")
}

var scenarioFile string
var externalServicesFile string

var (
  ErrLoadingRules = errors.New("unable to read rules file")
)

var scenarioCmd = &cobra.Command{
  Use:   "scenario",
  Short: "Run and evaluate a submt scenario",
  Long:  `Evaluate a submit scenario`,
  RunE: func(cmd *cobra.Command, args []string) error {
    studyRules, err := loadStudyRules()
    if(err != nil) {
      return errors.Join(ErrLoadingRules, err)
    }

    externalServices, err := evaluator.ReadExternalServicesFromYaml(externalServicesFile)
      
    if(err != nil) {
      return errors.Join(fmt.Errorf("unable to read '%s'", externalServicesFile), err)
    }
  
    scenarios, err := scenario.Load(scenarioFile)
    if(err != nil) {
      return errors.Join(fmt.Errorf("error loading scenarios in '%s'", scenarioFile), err)
    }
    
    for idx, sc := range scenarios {
      result := sc.Run(studyRules, externalServices)
      fmt.Printf("Scenario %d %s\n", idx, sc.Label)
      sc.PrintResult(result)  
    }
    return nil
  },
}
