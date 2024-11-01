package rules

import(
	"fmt"
	"github.com/influenzanet/study-service/pkg/types"
)

type RuleContext struct {
	EventType string
	Surveys []*RuleSurveyContext
	Current *RuleSurveyContext
	Default *RuleSurveyContext
}

type RuleSurveyContext struct {
	Survey string
	Keys map[string]struct{}
}

func NewSurveyContext(key string) *RuleSurveyContext {
	return &RuleSurveyContext{Survey: key, Keys: make(map[string]struct{}, 0) }
}

func FindRuleContext(exp types.Expression) *RuleContext {
	ctx := &RuleContext{Surveys: make([]*RuleSurveyContext, 0), Current: nil, Default: NewSurveyContext("_")}
	ctx.find(exp)
	return ctx
}

func isResponseExpression(name string) bool {
	switch(name) {
		case "responseHasKeysAny":
			return true
		case "responseHasOnlyKeysOtherThan":
			return true
		case "getResponseValueAsStr":
			return true
		case "getSelectedKeys":
			return true
		case "countResponseItems":
			return true
		case "hasResponseKey":
			return true
		case "hasResponseKeyWithValue":
			return true
		default:
			return false
	}
}

func (ctx *RuleContext) find(exp types.Expression) {
	for idx, data := range exp.Data {
		dtype := data.DType
		if(dtype == "") {
			dtype = "str"
		}
		if(dtype == "str") {
			strValue := data.Str
			if(exp.Name == "checkEventType" && idx == 0) {
				if(ctx.EventType == "") {
					ctx.EventType = strValue
				}
			}
			if(exp.Name == "checkSurveyResponseKey" && idx == 0) {
				surveyContext := NewSurveyContext(strValue)
				ctx.Surveys = append(ctx.Surveys, surveyContext)
				ctx.Current = surveyContext
			}
			if(isResponseExpression(exp.Name) && idx == 0) {
				key := strValue
				if(ctx.Current != nil) {
					ctx.Current.Keys[key] = struct{}{}
				} else {
					ctx.Default.Keys[key] = struct{}{}
				}		
			}
		}
		if(dtype == "exp") {
			ctx.find(*data.Exp)
		}
	}
}

func (sc *RuleSurveyContext) String() string {
	return fmt.Sprintf("%s: %s", sc.Survey, sc.Keys)
}