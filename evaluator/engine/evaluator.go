package engine

import (
	"fmt"
	"github.com/coneno/logger"
	"log"
	"bytes"
	"github.com/influenzanet/study-service/pkg/studyengine"
	"github.com/influenzanet/study-service/pkg/types"
	//pkgRules "github.com/grippenet/study-rules-evaluator/evaluator/rules"
)

type RuleEvaluator struct {
	Data []types.SurveyResponse
	Rules []types.Expression
	Verbose bool
}


func NewRuleEvaluator(previousData []types.SurveyResponse, rules []types.Expression) *RuleEvaluator {
	if(previousData == nil) {
		previousData = make([]types.SurveyResponse, 0, 0)
	}
	return &RuleEvaluator{Data: previousData, Rules: rules}
}


func (ev *RuleEvaluator) Submit(initialState types.ParticipantState, response types.SurveyResponse) EvaluationResult {
	instanceID := "dummy"
	studyKey := "dummy"
	
	dbService := MemoryDBService{Data: ev.Data}
	
	event := types.StudyEvent{
		InstanceID:                            instanceID,
		StudyKey:                              studyKey,
		Response: response,
		Type: "SUBMIT",
		ParticipantIDForConfidentialResponses: "",
	}

	actionData := studyengine.ActionData{
		PState:         initialState,
		ReportsToCreate: map[string]types.Report{},
	}

	actionConfig := studyengine.ActionConfigs{
		DBService:            dbService ,
		ExternalServiceConfigs: nil,
	}

	results := make([]EvalResult, 0)
	hasError := false

	oldLogger := logger.Debug
	logger.SetLevel(logger.LEVEL_DEBUG)

	for index, rule := range ev.Rules {

		var logBuffer bytes.Buffer
		l := log.New(&logBuffer, "", 0)

		logger.Debug = l

		if(ev.Verbose) {
			fmt.Printf("  Evaluate rule %d\n", index)
		}
		
		newState, err := studyengine.ActionEval(rule, actionData, event, actionConfig)

		r := EvalResult{
			Index: index,
			Data: CloneActionData(newState),
		}

		if(err != nil) {
			r.Error =  fmt.Sprintf("%s", err)
			hasError = true
		}

		r.Debug = logBuffer.String()

		if(r.Debug != "") {
			fmt.Println(r.Debug)
		}

		results = append(results, r )

		actionData.PState = newState.PState
	}

	logger.Debug = oldLogger

	return EvaluationResult{States: results, HasError: hasError }
}
