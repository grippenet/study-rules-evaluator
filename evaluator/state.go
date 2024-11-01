package evaluator

import(
	"time"
	"github.com/influenzanet/study-service/pkg/types"
)

func CreateEmptyParticipantState(flags map[string]string) types.ParticipantState {
	if(flags == nil) {
		flags = make(map[string]string, 0)
	}
	return types.ParticipantState{
		ParticipantID: "dummy",
		EnteredAt: time.Now().Unix(),
		StudyStatus: "active",
		Flags: flags,
		AssignedSurveys: make([]types.AssignedSurvey, 0),
		LastSubmissions: make(map[string]int64, 0),
		Messages: make([]types.ParticipantMessage, 0),
	}
} 


type StateBuilder struct {
	state *types.ParticipantState
}

func NewStateBuilder() *StateBuilder {
	part := CreateEmptyParticipantState(nil)
	return &StateBuilder{state: &part}
}

func (s *StateBuilder) AddFlag(name string, value string) {
	s.state.Flags[name] = value
}

func (s *StateBuilder) AddSurvey(surveyKey string, mode string) {
	as := types.AssignedSurvey{SurveyKey: surveyKey, Category: mode}
	s.state.AssignedSurveys = append(s.state.AssignedSurveys, as)
}

func (s *StateBuilder) State() types.ParticipantState {
	return *s.state
}