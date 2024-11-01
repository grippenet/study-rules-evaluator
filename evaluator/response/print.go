
package response

import(
	"github.com/influenzanet/study-service/pkg/types"
	"fmt"
	"strings"
)

// Print a Survey Response
func Print(sr types.SurveyResponse) {
	fmt.Println(sr.Key)
	fmt.Println("Responses:")
	for _, item := range sr.Responses {
		printSurveyItemResponse(item, 1)
	}
}

func printSurveyItemResponse(r types.SurveyItemResponse, level int) {
	fmt.Printf("%s- (SIR) '%s'\n", strings.Repeat(" ", level), r.Key)
	if(r.Response != nil) {
		printResponseItem(r.Response, level + 1)
	} else {
		for _, item := range r.Items {
			printSurveyItemResponse(item, level + 1)
		}
	}
}

func printResponseItem(r *types.ResponseItem, level int) {
	fmt.Printf("%s- (RI) '%s'", strings.Repeat(" ", level), r.Key)
	if(len(r.Items) > 0) {
		fmt.Println("")
		for _, item := range r.Items {
			printResponseItem(item, level + 1)
		}
	} else {
		fmt.Printf(" Dtype: '%s' Value: '%s'\n", r.Dtype, r.Value)
	}	
}

