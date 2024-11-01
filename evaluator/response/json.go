package response

import(
	"github.com/influenzanet/study-service/pkg/types"
	//"fmt"
	"os"
	"encoding/json"
)

type JsonSurveyResponse struct {
	Key           string               `bson:"key" json:"key"`
	Responses     []JsonSurveyItemResponse `bson:"responses" json:"responses"`
}

type JsonSurveyItemResponse struct {
	Key  string       `bson:"key" json:"key"`
	
	// for groups:
	Items []JsonSurveyItemResponse `bson:"items,omitempty" json:"items,omitempty"`

	// for single items:
	Response         *JsonResponseItem `bson:"response,omitempty" json:"response,omitempty"`
}

type JsonResponseItem struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value,omitempty" json:"value,omitempty"`
	Dtype string `bson:"dtype,omitempty" json:"dtype,omitempty"`
	// for response option groups
	Items []*JsonResponseItem `bson:"items,omitempty" json:"items,omitempty"`
}

func ReadSurveyResponseFromJSON(filename string) (JsonSurveyResponse, error) {
	content, err := os.ReadFile(filename)
	var input JsonSurveyResponse
	if err != nil {
		return input, err
	}
	err = json.Unmarshal(content, &input)
	return input, err
}

func ToSurveyResponse(from JsonSurveyResponse) *types.SurveyResponse {
	
	rr := make([]types.SurveyItemResponse, 0, len(from.Responses))
	
	for _, r := range from.Responses {
		rr = append(rr, toSurveyItemResponse(r))
	}

	sr := types.SurveyResponse{
		Key: from.Key,
		Responses: rr,
	}
	//fmt.Println("SurveyResponse ", sr)

	return &sr
}


func toSurveyItemResponse(from JsonSurveyItemResponse) types.SurveyItemResponse {
	r := types.SurveyItemResponse{
		Key: from.Key,
	}
	if(from.Response != nil) {
		r.Response = toResponseItem(from.Response)
	}
	if(len(from.Items) > 0) {
		rr := make([]types.SurveyItemResponse, 0, len(from.Items))
		for _, item := range from.Items {
			rr = append(rr, toSurveyItemResponse(item))
		}
		r.Items = rr
	}
	return r
}

func toResponseItem(from *JsonResponseItem) *types.ResponseItem {
	r := types.ResponseItem{
		Key: from.Key,	
	}
	if(len(from.Items) > 0) {
		rr := make([]*types.ResponseItem, 0, len(from.Items))
		for _, item := range from.Items {
			rr = append(rr, toResponseItem(item))
		}
		r.Items = rr
	}
	return &r
}
