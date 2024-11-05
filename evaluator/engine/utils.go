package engine

import (
	"github.com/influenzanet/study-service/pkg/studyengine"
	"github.com/influenzanet/study-service/pkg/types"
)

func CloneActionData(ad studyengine.ActionData) studyengine.ActionData {
	d := studyengine.ActionData{
		PState: CloneParticipantState(ad.PState),
		ReportsToCreate: cloneMap[string, types.Report](ad.ReportsToCreate),
	}
	return d
}

// CloneParticipantState make a deep copy of participant state
func CloneParticipantState(state types.ParticipantState) types.ParticipantState {
	newState := state
	newState.Flags = cloneMap[string,string](state.Flags)
	newState.LastSubmissions = cloneMap[string, int64](state.LastSubmissions)
	return newState
}

// cloneMap create a deep clone of a map
func cloneMap[K comparable, V any](m map[K]V) map[K]V {
	n := make(map[K]V, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}