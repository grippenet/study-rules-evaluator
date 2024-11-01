package engine

import (
	"github.com/influenzanet/study-service/pkg/studyengine"
	//"github.com/influenzanet/study-service/pkg/types"
)


type EvalResult struct {
	Index int `json:"index"`
	Data studyengine.ActionData`json:"result"`
	Error string `json:"error"`
	Debug string
}

type EvaluationResult struct {
	States []EvalResult
	HasError bool
}

func (ev *EvaluationResult) Last() *EvalResult {
	if(len(ev.States) == 0) {
		return nil
	}
	return  &ev.States[len(ev.States)-1]
}