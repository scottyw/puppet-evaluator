package types

import (
	"github.com/lyraproj/puppet-evaluator/eval"
)

func toFloat(v eval.Value) (float64, bool) {
	if iv, ok := v.(*FloatValue); ok {
		return iv.Float(), true
	}
	return 0.0, false
}

func toInt(v eval.Value) (int64, bool) {
	if iv, ok := v.(*IntegerValue); ok {
		return iv.Int(), true
	}
	return 0, false
}

func init() {
	eval.ToInt = toInt
	eval.ToFloat = toFloat
}
