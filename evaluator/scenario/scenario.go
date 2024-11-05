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
	"github.com/expr-lang/expr"
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
		if(sb.Data != nil) {
			sb.Response = response.ToSurveyResponse(*sb.Data)
		}
		if(sb.File != "") {
			fn := filepath.Join(dir, sb.File)
			data, err := response.ReadSurveyResponseFromJSON(fn)
			if err != nil {
				return errors.Join(fmt.Errorf("Error loading file in submit %d", i), err)
			}
			sb.Response = response.ToSurveyResponse(data)
		}
		if(sb.Response == nil) {
			return fmt.Errorf("No response provided for submit %d", i)
		}
	}
	err := scenario.Init()
	return err
}

func createAssertionEnv(state types.ParticipantState, previousState types.ParticipantState) map[string]any {
	return map[string]any {
		"previousState": previousState,
		"state": state,
	}
}

func evalAssertion(assertion string, env map[string]any) (bool, error) {
	program, err := expr.Compile(assertion, expr.Env(env), expr.AsBool())
	if err != nil {
		return false, err
	}
	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	return output.(bool), nil
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

func (sc *Scenario) Run(evaluator *engine.RuleEvaluator) *ScenarioResult {
	state := sc.State
	result := &ScenarioResult{Count: len(sc.Submits), Submits: make([]ScenarioSubmitResult, 0, len(sc.Submits))}

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
		if(!submitError) {
			last := evalResult.Last()

			//fmt.Println("Eval");
			//fmt.Printf("%+v\n", last.Data.PState.Flags);

			flagsChanges := change.CompareMap(state.Flags, last.Data.PState.Flags)

			submitResult.FlagsChanges = flagsChanges

			state = last.Data.PState

			env := createAssertionEnv(state, previousState)

			submitResult.Asserts = make([]AssertionResult, 0, len(submit.Assertions))

			for _, expectation := range submit.Assertions {
				b, err := evalAssertion(expectation, env)
				submitResult.Asserts = append(submitResult.Asserts, AssertionResult{Ok: b, Error: errAsString(err), } )
			}
			// Clone the state to be sure we keep map values as is
			rState := engine.CloneParticipantState(state)
			submitResult.State = &rState
		}
		result.HasError = result.HasError || len(submitResult.Errors) > 0 
		result.Submits = append(result.Submits, submitResult)
	}

	return result
}

func (sc *Scenario) PrintResult(r *ScenarioResult) {
	fmt.Printf("Scenario submits %d / %d\n", len(r.Submits), r.Count)
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
			printFlagsWithChange(submit.State.Flags, submit.FlagsChanges, 4)
		}
		if(len(submit.Asserts) > 0) {
			fmt.Println("   Assertions")
			printAssertions(submitDef.Assertions, submit.Asserts, 4)
		}
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