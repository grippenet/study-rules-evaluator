package scenario

import(
	"time"
	"github.com/influenzanet/study-service/pkg/types"
	"github.com/grippenet/study-rules-evaluator/evaluator/response"
)

// TimeRefSpec hold time reference in scenario
// Time holds an absolute value point of time, used if defined
// Relative a time point relative to an external time (usually previous time point)
// Duration holds shift using go time.Duration (only hours are handled)
// Days & weeks are just shortcuts to add more convenient time unit (total time shift is duration + days + weeks)
type TimeRefSpec struct {
	Fixed string `json:"fixed"` // Absolute time as ISO Time
	Duration string `json:"duration"` // Add duration using Duration syntax
	Days int `json:"days"`
	Weeks int `json:"weeks"`
}

type SubmitResponse struct {
	Data *response.JsonSurveyResponse `json:"data"`
	File string `json:"file"`
	TimeSpec *TimeRefSpec `json:"time"`
	FillingDuration int `json:"filling"` // Number of seconds spends by the user to fill & send the survey
	Response *types.SurveyResponse
	Assertions []string `json:"asserts"`
	timeRef TimeRef // parsed TimeReference
}

type Scenario struct {
	Time string  `json:"time"` // The base time of the scenario 
	State  types.ParticipantState `json:"state"`
	Submits []SubmitResponse `json:"submits"`
	startTime time.Time // parsed Time
}

type TimeRef interface {
	ShiftTime(previous time.Time) time.Time
}

type ScenarioResult struct {
	Count int
	HasError bool
	Submits []ScenarioSubmitResult
}

type AssertionResult struct {
	Ok bool
	Error string
}

type ScenarioSubmitResult struct {
	Submit int  `json:"submit"`
	Time   time.Time `json:"time"`
	Errors  []ScenarioError `json:"errors"`
	Asserts []AssertionResult `json:"assertions"`
	State *types.ParticipantState // Final state of participant
	FlagsChanges map[string]int
}

type ScenarioError struct {
	Submit int `json:"submit"`
	Type string `json:"type"`
	Index int `json:"index"`
	Error string `json:"error"`
}