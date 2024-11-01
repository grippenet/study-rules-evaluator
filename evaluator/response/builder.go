package response

import(
	"github.com/influenzanet/study-service/pkg/types"
)

type Builder struct {

}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SurveyResponse(surveyKey string, time int64, responses ...types.SurveyItemResponse) *types.SurveyResponse {
	return &types.SurveyResponse{
		Key: surveyKey,
		ParticipantID: "dummy",
		VersionID:"01-00",
		OpenedAt: time - 60,
		SubmittedAt: time,
		ArrivedAt: time,
		Responses: responses,
	}
}

// SurveySingleItem
func (b *Builder) I(key string, response *types.ResponseItem) types.SurveyItemResponse {
	return types.SurveyItemResponse{
		Key: key,
		Response: response,
	}
}

// SurveyItemGroup
func (b *Builder) IG(key string, items ...types.SurveyItemResponse) types.SurveyItemResponse {
	return types.SurveyItemResponse{
		Key: key,
		Items: items,
	}
}

// ResponseItem with only a key
func (b *Builder) R(key string) *types.ResponseItem {
	return &types.ResponseItem{
		Key: key,
	}
}

// ResponseItem with Value
func (b *Builder) V(key string, value string, dtype string) *types.ResponseItem {
	return &types.ResponseItem{
		Key: key,
		Value: value,
		Dtype: dtype,
	}
}

// ResponseItem with Value
func (b *Builder) RG(key string, items ...*types.ResponseItem) *types.ResponseItem {
	return &types.ResponseItem{
		Key: key,
		Items: items,
	}
}
