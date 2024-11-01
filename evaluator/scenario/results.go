package scenario

import(

	"github.com/influenzanet/study-service/pkg/types"
)

type ScenarioResult struct {
	Count int
	HasError bool
	Submits []ScenarioSubmitResult
}

type ExpectationResult struct {
	Ok bool
	Error string
}

type ScenarioSubmitResult struct {
	Submit int  `json:"submit"`
	Errors  []ScenarioError `json:"errors"`
	Expects []ExpectationResult `json:"expectations"`
	State *types.ParticipantState // Final state of participant
	FlagsChanges map[string]int
}

type ScenarioError struct {
	Submit int `json:"submit"`
	Type string `json:"type"`
	Index int `json:"index"`
	Error string `json:"error"`
}