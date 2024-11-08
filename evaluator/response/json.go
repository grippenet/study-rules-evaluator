package response

import(
	"github.com/influenzanet/study-service/pkg/types"
	//"fmt"
	"os"
	"encoding/json"
)

func ReadSurveyResponseFromJSON(filename string) (types.SurveyResponse, error) {
	content, err := os.ReadFile(filename)
	var input types.SurveyResponse
	if err != nil {
		return input, err
	}
	err = json.Unmarshal(content, &input)
	return input, err
}