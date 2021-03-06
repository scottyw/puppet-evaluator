package types

import (
	"io"

	"github.com/lyraproj/puppet-evaluator/eval"
)

type ScalarDataType struct{}

var ScalarData_Type eval.ObjectType

func init() {
	ScalarData_Type = newObjectType(`Pcore::ScalarDataType`, `Pcore::ScalarType{}`, func(ctx eval.Context, args []eval.Value) eval.Value {
		return DefaultScalarDataType()
	})
}

func DefaultScalarDataType() *ScalarDataType {
	return scalarDataType_DEFAULT
}

func (t *ScalarDataType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
}

func (t *ScalarDataType) Equals(o interface{}, g eval.Guard) bool {
	_, ok := o.(*ScalarDataType)
	return ok
}

func (t *ScalarDataType) IsAssignable(o eval.Type, g eval.Guard) bool {
	switch o.(type) {
	case *ScalarDataType:
		return true
	default:
		return GuardedIsAssignable(stringType_DEFAULT, o, g) ||
			GuardedIsAssignable(integerType_DEFAULT, o, g) ||
			GuardedIsAssignable(booleanType_DEFAULT, o, g) ||
			GuardedIsAssignable(floatType_DEFAULT, o, g)
	}
}

func (t *ScalarDataType) IsInstance(o eval.Value, g eval.Guard) bool {
	switch o.(type) {
	case *BooleanValue, *FloatValue, *IntegerValue, *StringValue:
		return true
	}
	return false
}

func (t *ScalarDataType) MetaType() eval.ObjectType {
	return ScalarData_Type
}

func (t *ScalarDataType) Name() string {
	return `ScalarData`
}

func (t *ScalarDataType)  CanSerializeAsString() bool {
  return true
}

func (t *ScalarDataType)  SerializationString() string {
	return t.String()
}


func (t *ScalarDataType) String() string {
	return `ScalarData`
}

func (t *ScalarDataType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *ScalarDataType) PType() eval.Type {
	return &TypeType{t}
}

var scalarDataType_DEFAULT = &ScalarDataType{}
