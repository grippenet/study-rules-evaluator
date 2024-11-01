package engine

import(
	"github.com/influenzanet/study-service/pkg/dbs/studydb"
	"github.com/influenzanet/study-service/pkg/types"
	"sort"
)

type MemoryDBService struct {
	Data []types.SurveyResponse
}

func (m MemoryDBService) FindSurveyResponses(instanceID string, studyKey string, query studydb.ResponseQuery) (responses []types.SurveyResponse, err error) {
	selecteData := make([]types.SurveyResponse, 0)

	for _, r := range m.Data {
		if query.ParticipantID != "" && query.ParticipantID != r.ParticipantID {
			continue
		}
		if(query.SurveyKey != "" && r.Key != query.SurveyKey) {
			continue
		}
		submittedAt := r.SubmittedAt
		keep := true
		if query.Since > 0 && query.Until > 0 {
			keep = submittedAt > query.Since && submittedAt < query.Until
			// filter["$and"] = bson.A{
			//	bson.M{"submittedAt": bson.M{"$gt": query.Since}},
			//	bson.M{"submittedAt": bson.M{"$lt": query.Until}},
			// }
		} else if query.Since > 0 {
			keep = submittedAt > query.Since
			// filter["submittedAt"] = bson.M{"$gt": query.Since}
		} else if query.Until > 0 {
			keep = submittedAt < query.Until
			// filter["submittedAt"] = bson.M{"$lt": query.Until}
		}
		if(!keep) {
			continue
		}
		selecteData = append(selecteData, r)
	}

	// Sort in reverse order of submission time
	sort.SliceStable(selecteData, func(i,j int) bool {
		return selecteData[i].SubmittedAt > selecteData[i].SubmittedAt
	})

	n := len(selecteData)
	if(query.Limit > 0 && n >= int(query.Limit)) {
		n = int(query.Limit)
		selecteData = selecteData[:n]
	}
	responses = selecteData

	return responses, nil
}

func (m MemoryDBService) DeleteConfidentialResponses(instanceID string, studyKey string, participantID string, key string) (count int64, err error) {
	return 0, nil
}

func (m MemoryDBService) SaveResearcherMessage(instanceID string, studyKey string, message types.StudyMessage) error {
	return nil
}
