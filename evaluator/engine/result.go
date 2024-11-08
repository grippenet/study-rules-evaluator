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

type DebugMessage struct {
	Index int `json:"index"`
	Message string `json:"msg"`
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

func (ev *EvaluationResult) CollectDebug() []DebugMessage {
	msgs := make([]DebugMessage, 0)
	if(len(ev.States) == 0) {
		return nil
	}
	for _, s := range ev.States {
		if(s.Debug != "") {
			msgs = append(msgs, DebugMessage{Index: s.Index, Message: s.Debug})
		}
	}
	return msgs
}
