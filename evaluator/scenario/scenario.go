package scenario

import(
	"os"
	"fmt"
	"time"
	"strings"
	"errors"
	"encoding/json"
	"path/filepath"
	"github.com/fatih/color"
	"github.com/influenzanet/study-service/pkg/types"
	"github.com/influenzanet/study-service/pkg/studyengine"
	"github.com/grippenet/study-rules-evaluator/evaluator/engine"
	"github.com/grippenet/study-rules-evaluator/evaluator/response"
	"github.com/grippenet/study-rules-evaluator/evaluator/change"
)

func readScenarioFromJSON(file string) ([]Scenario, error) {
	content, err := os.ReadFile(file)
	var input []Scenario
	if err != nil {
		return input, err
	}
	err = json.Unmarshal(content, &input)
	return input, err
}

func Load(file string) ([]Scenario, error) {
	scenarios, err := readScenarioFromJSON(file)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(file)
	for idx, _ := range scenarios {
		sc := &scenarios[idx]
		err = initScenario(sc, dir)
		if(err != nil) {
			return nil, errors.Join( fmt.Errorf("Error in entry %d of '%s'", idx, file), err)
		}
	}
	return scenarios, err
}

func initScenario(scenario *Scenario, dir string) error {
	for i, _ := range scenario.Submits {
		sb := &scenario.Submits[i]
		if(sb.File != "") {
			if(sb.Response != nil) {
				return errors.New("Cannot provide `file` and `data` field, cannot choice")
			}
			fn := filepath.Join(dir, sb.File)
			data, err := response.ReadSurveyResponseFromJSON(fn)
			if err != nil {
				return errors.Join(fmt.Errorf("Error loading file in submit %d", i), err)
			}
			sb.Response = &data
		}
		if(sb.Response == nil) {
			return fmt.Errorf("No response provided for submit %d", i)
		}
	}
	err := scenario.Init()
	return err
}

// Init parse & performs checks before the scenario is run
func (sc *Scenario) Init() error {
	var err error
	if(sc.Time == "") {
		sc.startTime = time.Now()
	} else {
		sc.startTime, err = time.Parse(FixedTimeLayout, sc.Time)
		if err != nil {
			return errors.Join(errors.New("Unable to parse scenario.Time field"), err)
		}
	}
	// Check submit fields are parseable
	for idx, _ := range sc.Submits {
		sb := &sc.Submits[idx]
		err := sb.Init()
		if err != nil {
			return errors.Join(fmt.Errorf("Error in submit %d", idx), err)
		}
	}
	return nil
}

func (sc *Scenario) SetVerbose(v bool) {
	sc.verbose = v
} 

func (sc *Scenario) Run(studyRules []types.Expression, ExternalServiceConfigs []types.ExternalService) *ScenarioResult {
	state := sc.State
	result := &ScenarioResult{Count: len(sc.Submits), Submits: make([]ScenarioSubmitResult, 0, len(sc.Submits))}

	dbService := engine.NewMemoryDBService()

	evaluator := engine.NewRuleEvaluator(dbService, studyRules)
	if(len(ExternalServiceConfigs) > 0) {
		evaluator.WithExternalServices(ExternalServiceConfigs)
	}
	//evaluator.Verbose = true
	
	now := sc.startTime
	for idx, _ := range sc.Submits {

		submit := &sc.Submits[idx]

		now = submit.ShiftTime(now)

		// Change study engine current time
		studyengine.Now = func() time.Time {
			return now
		}

		submitResult := ScenarioSubmitResult{Submit: idx, Errors: make([]ScenarioError, 0), Time: now, }
		submitError := false
		
		fmt.Printf("= Submit %d '%s' at %s\n", idx, submit.Response.Key, now)
		
		if(sc.verbose) {
			response.Print(*submit.Response)
		}

		u := now.Unix()
		submit.Response.OpenedAt = u - int64(submit.FillingDuration)
		submit.Response.SubmittedAt = u
		submit.Response.ArrivedAt = u

		previousState := state
		
		evalResult := evaluator.Submit(state, *submit.Response)

		if(evalResult.HasError) {
			fmt.Println("Errors found")
			for _, ss := range evalResult.States {
				submitResult.Errors = append(submitResult.Errors, ScenarioError{Type: "rule", Index: ss.Index, Error: ss.Error})
			}
			submitError = true
		}

		submitResult.DebugMessages = evalResult.CollectDebug()

		if(!submitError) {
			last := evalResult.Last()

			//fmt.Println("Eval");
			//fmt.Printf("%+v\n", last.Data.PState.Flags);

			flagsChanges := change.CompareMap(state.Flags, last.Data.PState.Flags)

			submitResult.FlagsChanges = flagsChanges
			

			state = last.Data.PState

			assertionEnv := AssertionEnv{
				State: state, 
				PreviousState: previousState,
				Reports: last.Data.ReportsToCreate,
				SubmitAt: now,
			}
			
			submitResult.Asserts = make([]AssertionResult, 0, len(submit.Assertions))

			for _, expectation := range submit.Assertions {
				b, err := evalAssertion(expectation, assertionEnv)
				submitResult.Asserts = append(submitResult.Asserts, AssertionResult{Ok: b, Error: errAsString(err), } )
			}
			// Clone the state to be sure we keep map values as is
			rState := engine.CloneParticipantState(state)
			submitResult.State = &rState
			submitResult.Reports = last.Data.ReportsToCreate
		}
		result.HasError = result.HasError || len(submitResult.Errors) > 0 
		result.Submits = append(result.Submits, submitResult)
	}

	return result
}

func (sc *Scenario) PrintResult(r *ScenarioResult) {
	fmt.Printf("Scenario submits %d / %d\n", len(r.Submits), r.Count)
	indent := 4
	for idx, submit := range r.Submits {
		fmt.Printf(" > Submit %d at '%s' ==> \n", idx, submit.Time.Format("2006-02-01 15:04:05"))
		submitDef := sc.Submits[idx]
		if(len(submit.Errors) > 0) {
			fmt.Println("  Errors: ")
			for _, e := range submit.Errors {
				fmt.Printf(" - %s[%d] %s ", e.Type, e.Index, e.Error)
			}
		}
		if(submit.State != nil) {
			fmt.Println("  Flags")
			printFlagsWithChange(submit.State.Flags, submit.FlagsChanges, indent)
		}
		if(len(submit.Reports) > 0) {
			fmt.Println("  Reports")
			printReports(submit.Reports, indent)
			//fmt.Println(submit.Reports)		
		}
		if(len(submit.Asserts) > 0) {
			fmt.Println("  Assertions")
			printAssertions(submitDef.Assertions, submit.Asserts, indent)
		}
		if(len(submit.DebugMessages) > 0) {
			fmt.Println("  Debug messages")
			printDebugMessages(submit.DebugMessages, indent)
		}
	}
}

func reportToString(report types.Report ) string {
	var sb strings.Builder
	t := time.Unix(report.Timestamp, 0)
	sb.WriteString(t.Format(FixedTimeLayout))
	for _, d := range report.Data {
		var dtype string
		if(d.Dtype != "") {
			dtype = fmt.Sprintf("[%s]", d.Dtype)
		}
		sb.WriteString(fmt.Sprintf(" `%s`%s=%s", d.Key, dtype, d.Value))
	}
	return sb.String()
}

func printReports(reports map[string]types.Report, indent int) {
	prefix := strings.Repeat(" ", indent)
	for key, report := range reports {
		fmt.Printf("%s- `%s` : %s\n", prefix, key, reportToString(report))
	}
}

func printFlagsWithChange(flags map[string]string, changes map[string]int, indent int) {
	prefix := strings.Repeat(" ", indent)
	for _, name := range sortedMapKeys(changes) {
		mod, _ := changes[name]
		v, _ := flags[name]
		tag := "?"
		showValue := true
		c := color.New(color.FgBlack)
		switch(mod) {
			case change.Equal:
				tag = "="
			case change.Changed:
				tag = "*"
				c = color.New(color.FgYellow)
			case change.Added:
				tag = "+"
				c = color.New(color.FgGreen)
			case change.Deleted:
				tag = "-"
				showValue = false
				c = color.New(color.FgRed, color.CrossedOut)
		}
		if(showValue) {
			c.Printf("%s%s %s = %s\n", prefix, tag, name, v)
		} else {
			c.Printf("%s%s %s\n",prefix, tag, name)
		}
	}
}

func printAssertions(definitions []string, asserts []AssertionResult, indent int) {
	prefix := strings.Repeat(" ", indent)
	for idx, r := range asserts {
		def := definitions[idx]
		if(r.Error != "") {
			color.Red("%s- `%s` <error> %s)\n", prefix, def, r.Error)
		} else {
			var col color.Attribute
			if(r.Ok) {
				col = color.FgGreen
			} else {
				col = color.FgYellow
			}
			color.New(col).Printf("%s- `%s` = %t\n", prefix, def, r.Ok)
		}
	}
}

func printDebugMessages(messages []engine.DebugMessage, indent int) {
	prefix := strings.Repeat(" ", indent)
	for _, m := range messages {
		fmt.Printf("%s- Rule %d : %s\n", prefix, m.Index, m.Message)
	}
}