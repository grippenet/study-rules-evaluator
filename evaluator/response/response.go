package response

import(
	"github.com/influenzanet/study-service/pkg/types"
	//"errors"
	"strings"
)

type SurveyResponseBuilder struct {
	resp *types.SurveyResponse
}

func CreateEmptyResponse(surveyKey string, time int64) types.SurveyResponse {
	return types.SurveyResponse{
		Key: surveyKey,
		ParticipantID: "dummy",
		VersionID:"01-00",
		OpenedAt: time - 60,
		SubmittedAt: time,
		ArrivedAt: time,
		Responses: make([]types.SurveyItemResponse, 0),
	}
}

func NewResponseBuilder(surveyKey string, time int64) *SurveyResponseBuilder {
	s := CreateEmptyResponse(surveyKey, time)
	return &SurveyResponseBuilder{resp: &s}
}

func (r *SurveyResponseBuilder) Response() *types.SurveyResponse {
	return r.resp
}

func (r * SurveyResponseBuilder) Group(key string) *SurveyGroupItemResponseBuilder {
	itemResponse := types.SurveyItemResponse{Key: key}
	r.resp.Responses = append(r.resp.Responses, itemResponse)
	return NewSurveyGroupItemResponseBuilder(&r.resp.Responses[len(r.resp.Responses)-1])
}

func (r * SurveyResponseBuilder) Single(key string) *SurveySingleItemResponseBuilder {
	itemResponse := types.SurveyItemResponse{Key: key}
	r.resp.Responses = append(r.resp.Responses, itemResponse)
	return NewSurveySingleItemResponseBuilder(&r.resp.Responses[len(r.resp.Responses)-1])
}

func (r *SurveyResponseBuilder) AddSingleItem(itemKey string, groupsKeys string, responseKeys ...string) *types.ResponseItem {
	item := types.SurveyItemResponse{
		Key: itemKey,
	}
	root := &types.ResponseItem{}
	item.Response = root
	groups := strings.Split(groupsKeys, ".")

	for _, groupKey := range groups {
		root.Key = groupKey
		next := &types.ResponseItem{}
		responseItems := make([]*types.ResponseItem, 0, 0)
		responseItems = append(responseItems, next)
		root.Items = responseItems
		root = next
	}

	root.Items = make([]*types.ResponseItem, 0, len(responseKeys))
	for _, responseKey := range responseKeys {
		ri := &types.ResponseItem{Key: responseKey}
		root.Items = append(root.Items, ri)
	}
	r.resp.Responses = append(r.resp.Responses, item)
	return root
}

type SurveyGroupItemResponseBuilder struct {
	resp *types.SurveyItemResponse
}

func NewSurveyGroupItemResponseBuilder(resp *types.SurveyItemResponse) *SurveyGroupItemResponseBuilder {
	return &SurveyGroupItemResponseBuilder{resp: resp}
}

func (b *SurveyGroupItemResponseBuilder) Group(key string) *SurveyGroupItemResponseBuilder {
	itemResponse := types.SurveyItemResponse{Key: key}
	b.resp.Items = append(b.resp.Items, itemResponse)
	return NewSurveyGroupItemResponseBuilder(&b.resp.Items[len(b.resp.Items)-1])
}

func (b *SurveyGroupItemResponseBuilder) Single(key string) *SurveySingleItemResponseBuilder {
	itemResponse := types.SurveyItemResponse{Key: key}
	b.resp.Items = append(b.resp.Items, itemResponse)
	return NewSurveySingleItemResponseBuilder(&b.resp.Items[len(b.resp.Items)-1])
}

type SurveySingleItemResponseBuilder struct {
	resp *types.SurveyItemResponse
}

func NewSurveySingleItemResponseBuilder(resp *types.SurveyItemResponse) *SurveySingleItemResponseBuilder {
	return &SurveySingleItemResponseBuilder{resp: resp}
}

func (b *SurveySingleItemResponseBuilder) ResponseSingle(key string) *ResponseItemSingleBuilder {
	rb := NewResponseItemSingleBuilder(key)
	b.resp.Response = rb.resp
	return rb
}

func (b *SurveySingleItemResponseBuilder) ResponseGroup(key string) *ResponseItemGroupBuilder {
	rb := NewResponseItemGroupBuilder(key)
	b.resp.Response = rb.resp
	return rb
}


type ResponseItemGroupBuilder struct {
	resp *types.ResponseItem
}

func NewResponseItemGroupBuilder(key string) *ResponseItemGroupBuilder {
	resp := types.ResponseItem{Key: key}
	return &ResponseItemGroupBuilder{resp: &resp}
}

func (b *ResponseItemGroupBuilder) Group(key string) *ResponseItemGroupBuilder {
	child := NewResponseItemGroupBuilder(key)
	b.resp.Items = append(b.resp.Items , child.resp) 
	return child
}

func (b *ResponseItemGroupBuilder) Single(key string) *ResponseItemSingleBuilder {
	child := NewResponseItemSingleBuilder(key)
	b.resp.Items = append(b.resp.Items, child.resp) 
	return child
}

type ResponseItemSingleBuilder struct {
	resp *types.ResponseItem
}

func NewResponseItemSingleBuilder(key string) *ResponseItemSingleBuilder {
	resp := types.ResponseItem{Key: key}
	return &ResponseItemSingleBuilder{resp: &resp}
}

func (b *ResponseItemSingleBuilder) Value(v string, dt string) {
	b.resp.Value = v
	b.resp.Dtype = dt
}
