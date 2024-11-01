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
	"github.com/grippenet/study-rules-evaluator/evaluator/engine"
	"github.com/grippenet/study-rules-evaluator/evaluator/response"
	"github.com/grippenet/study-rules-evaluator/evaluator/change"
	"github.com/expr-lang/expr"
)

type SubmitResponse struct {
	Data *response.JsonSurveyResponse `json:"data"`
	File string `json:"file"`
	Time string
	Response *types.SurveyResponse
	Expectations []string `json:"expect"`
}

type Scenario struct {
	State  types.ParticipantState `json:"state"`
	Submits []SubmitResponse `json:"submits"` 	
}

func readScenarioFromJSON(file string) (Scenario, error) {
	content, err := os.ReadFile(file)
	var input Scenario
	if err != nil {
		return input, err
	}
	err = json.Unmarshal(content, &input)
	return input, err
}


func Load(file string) (*Scenario, error) {
	scenario, err := readScenarioFromJSON(file)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(file)
	for i, _ := range scenario.Submits {
		sb := &scenario.Submits[i]
		if(sb.Data != nil) {
			sb.Response = response.ToSurveyResponse(*sb.Data)
		}
		if(sb.File != "") {
			fn := filepath.Join(dir, sb.File)
			data, err := response.ReadSurveyResponseFromJSON(fn)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("Error loading file in submit %d", i), err)
			}
			sb.Response = response.ToSurveyResponse(data)
		}
		if(sb.Response == nil) {
			return nil, fmt.Errorf("No response provided for submit %d", i)
		}
		if sb.Time != "" {
			time, err := time.Parse("2006-01-02T15:04:05", sb.Time)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("Unable to parse time in submit %d", i), err)
			}
			sb.Response.SubmittedAt = time.Unix()
			sb.Response.ArrivedAt = time.Unix()
		} 
	}
	return &scenario, nil
}

func createExpectationEnv(flags map[string]string) map[string]any {
	return map[string]any {
		"flags": flags,
	}
}


func evalExpectation(expectation string, env map[string]any) (bool, error) {
	program, err := expr.Compile(expectation, expr.Env(env), expr.AsBool())
	if err != nil {
		return false, err
	}
	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	return output.(bool), nil
}

func (sc *Scenario) Run(evaluator *engine.RuleEvaluator) *ScenarioResult {
	state := sc.State
	result := &ScenarioResult{Count: len(sc.Submits), Submits: make([]ScenarioSubmitResult, 0, len(sc.Submits))}
	for idx, submit := range sc.Submits {
		submitResult := ScenarioSubmitResult{Submit: idx, Errors: make([]ScenarioError, 0), }
		submitError := false
		fmt.Printf("= Submit %d\n", idx)
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

			flagsChanges := change.CompareMap(state.Flags, last.Data.PState.Flags)

			submitResult.FlagsChanges = flagsChanges

			state = last.Data.PState

			env := createExpectationEnv(state.Flags)

			submitResult.Expects = make([]ExpectationResult, 0, len(submit.Expectations))

			for _, expectation := range submit.Expectations {
				b, err := evalExpectation(expectation, env)
				submitResult.Expects = append(submitResult.Expects, ExpectationResult{Ok: b, Error: errAsString(err), } )
			}
			submitResult.State = &state
		}
		result.HasError = result.HasError || len(submitResult.Errors) > 0 
		result.Submits = append(result.Submits, submitResult)
	}
	return result
}

func (sc *Scenario) PrintResult(r *ScenarioResult) {
	fmt.Printf("Scenario submits %d / %d\n", len(r.Submits), r.Count)
	for idx, submit := range r.Submits {
		fmt.Printf(" > Submit %d ==> \n", idx)
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
		if(len(submit.Expects) > 0) {
			fmt.Println("   Expectations")
			printExpectations(submitDef.Expectations, submit.Expects, 4)
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

func printExpectations(definitions []string, expects []ExpectationResult, indent int) {
	prefix := strings.Repeat(" ", indent)
	for idx, r := range expects {
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