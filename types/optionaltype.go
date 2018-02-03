package types

import (
	"io"

	"github.com/puppetlabs/go-evaluator/errors"
	"github.com/puppetlabs/go-evaluator/eval"
)

type OptionalType struct {
	typ eval.PType
}

func DefaultOptionalType() *OptionalType {
	return optionalType_DEFAULT
}

func NewOptionalType(containedType eval.PType) *OptionalType {
	if containedType == nil || containedType == anyType_DEFAULT {
		return DefaultOptionalType()
	}
	return &OptionalType{containedType}
}

func NewOptionalType2(args ...eval.PValue) *OptionalType {
	switch len(args) {
	case 0:
		return DefaultOptionalType()
	case 1:
		if containedType, ok := args[0].(eval.PType); ok {
			return NewOptionalType(containedType)
		}
		if containedType, ok := args[0].(*StringValue); ok {
			return NewOptionalType3(containedType.String())
		}
		panic(NewIllegalArgumentType2(`Optional[]`, 0, `Variant[Type,String]`, args[0]))
	default:
		panic(errors.NewIllegalArgumentCount(`Optional[]`, `0 - 1`, len(args)))
	}
}

func NewOptionalType3(str string) *OptionalType {
	return &OptionalType{NewStringType(nil, str)}
}

func (t *OptionalType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
	t.typ.Accept(v, g)
}

func (t *OptionalType) ContainedType() eval.PType {
	return t.typ
}

func (t *OptionalType) Default() eval.PType {
	return optionalType_DEFAULT
}

func (t *OptionalType) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*OptionalType); ok {
		return t.typ.Equals(ot.typ, g)
	}
	return false
}

func (t *OptionalType) Generic() eval.PType {
	return NewOptionalType(eval.GenericType(t.typ))
}

func (t *OptionalType) IsAssignable(o eval.PType, g eval.Guard) bool {
	return GuardedIsAssignable(o, undefType_DEFAULT, g) || GuardedIsAssignable(t.typ, o, g)
}

func (t *OptionalType) IsInstance(o eval.PValue, g eval.Guard) bool {
	return o == _UNDEF || GuardedIsInstance(t.typ, o, g)
}

func (t *OptionalType) Name() string {
	return `Optional`
}

func (t *OptionalType) Parameters() []eval.PValue {
	if t.typ == DefaultAnyType() {
		return eval.EMPTY_VALUES
	}
	if str, ok := t.typ.(*StringType); ok && str.value != `` {
		return []eval.PValue{WrapString(str.value)}
	}
	return []eval.PValue{t.typ}
}

func (t *OptionalType) String() string {
	return eval.ToString2(t, NONE)
}

func (t *OptionalType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *OptionalType) Type() eval.PType {
	return &TypeType{t}
}

var optionalType_DEFAULT = &OptionalType{typ: anyType_DEFAULT}
