package evaluator

import(
	//"github.com/coneno/logger"
	
	"github.com/influenzanet/study-service/pkg/types"
	"os"
	"encoding/json"
	"gopkg.in/yaml.v2"
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

// Read the External service definition from yaml file (study-service format)
func ReadExternalServicesFromYaml(filename string) ([]types.ExternalService, error) {
	if(filename == "") {
		return []types.ExternalService{}, nil
	}
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var services types.ExternalServices
	err = yaml.UnmarshalStrict(yamlFile, &services)
	if err != nil {
		return nil, err
	}
	return services.Services, nil
}