package rules

import(
	"strings"
	"fmt"
	"github.com/influenzanet/study-service/pkg/types"
)

func ExpressionToString(expr types.Expression) string {
	var sb strings.Builder
	sb.WriteString(expr.Name)
	sb.WriteString("(")
	last := len(expr.Data) - 1
	for i, arg := range expr.Data {
		sb.WriteString(ExpressionArgToString(arg))
		if(i < last) {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func ExpressionArgToString(arg types.ExpressionArg) string {
	if arg.DType == "" || arg.DType == "str" {
		return fmt.Sprintf("\"%s\"", arg.Str)
	}
	if arg.DType == "num" {
		return fmt.Sprintf("%f", arg.Num)
	}
	if(arg.DType == "exp") {
		return ExpressionToString(*arg.Exp)
	}
	return "?"
}