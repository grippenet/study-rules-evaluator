package scenario

import(
	"time"
	"github.com/influenzanet/study-service/pkg/types"
	"github.com/expr-lang/expr"
)

type AssertionEnv struct {
	State types.ParticipantState `expr:"state"`
	PreviousState types.ParticipantState `expr:"previousState"`
	Reports map[string]types.Report `expr:"reports"`
	SubmitAt time.Time `expr:"submitAt"`
}

func (env *AssertionEnv) HasReport(reportKey string, key string) bool {
	r, ok := env.Reports[reportKey]
	if(!ok) {
		return false
	}
	if(key == "") {
		return true
	}
	for _, d := range r.Data {
		if(d.Key == key) {
			return true
		}
	}
	return false
}


func evalAssertion(assertion string, data AssertionEnv) (bool, error) {
	
	program, err := expr.Compile(assertion, expr.Env(data), expr.AsBool())
	if err != nil {
		return false, err
	}
	output, err := expr.Run(program, data)
	if err != nil {
		return false, err
	}
	return output.(bool), nil
}
