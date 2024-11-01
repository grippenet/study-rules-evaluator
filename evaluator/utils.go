package evaluator

import(
	//"github.com/coneno/logger"
	
	"github.com/influenzanet/study-service/pkg/types"
	"os"
	"encoding/json"
)

func ReadRulesFromJSON(filename string) ([]types.Expression, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var input []types.Expression
	err = json.Unmarshal(content, &input)
	return input, err
}