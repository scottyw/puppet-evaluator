package types

import (
	"bytes"
	"io"

	"github.com/lyraproj/puppet-evaluator/errors"
	"github.com/lyraproj/puppet-evaluator/eval"
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/semver/semver"
	"reflect"
	"github.com/lyraproj/puppet-evaluator/utils"
)

type (
	SemVerType struct {
		vRange semver.VersionRange
	}

	SemVerValue SemVerType
)

var semVerType_DEFAULT = &SemVerType{semver.MatchAll}

var SemVer_Type eval.ObjectType

func init() {
	SemVer_Type = newObjectType(`Pcore::SemVerType`,
		`Pcore::ScalarType {
	attributes => {
		ranges => {
      type => Array[Variant[SemVerRange,String[1]]],
      value => []
    }
	}
}`, func(ctx eval.Context, args []eval.Value) eval.Value {
			return NewSemVerType2(args...)
		})

	newGoConstructor2(`SemVer`,
		func(t eval.LocalTypes) {
			t.Type(`PositiveInteger`, `Integer[0,default]`)
			t.Type(`SemVerQualifier`, `Pattern[/\A(?<part>[0-9A-Za-z-]+)(?:\.\g<part>)*\Z/]`)
			t.Type(`SemVerString`, `String[1]`)
			t.Type(`SemVerHash`, `Struct[major=>PositiveInteger,minor=>PositiveInteger,patch=>PositiveInteger,Optional[prerelease]=>SemVerQualifier,Optional[build]=>SemVerQualifier]`)
		},

		func(d eval.Dispatch) {
			d.Param(`SemVerString`)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				v, err := semver.ParseVersion(args[0].String())
				if err != nil {
					panic(errors.NewIllegalArgument(`SemVer`, 0, err.Error()))
				}
				return WrapSemVer(v)
			})
		},

		func(d eval.Dispatch) {
			d.Param(`PositiveInteger`)
			d.Param(`PositiveInteger`)
			d.Param(`PositiveInteger`)
			d.OptionalParam(`SemVerQualifier`)
			d.OptionalParam(`SemVerQualifier`)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				argc := len(args)
				major := args[0].(*IntegerValue).Int()
				minor := args[1].(*IntegerValue).Int()
				patch := args[2].(*IntegerValue).Int()
				preRelease := ``
				build := ``
				if argc > 3 {
					preRelease = args[3].String()
					if argc > 4 {
						build = args[4].String()
					}
				}
				v, err := semver.NewVersion3(int(major), int(minor), int(patch), preRelease, build)
				if err != nil {
					panic(errors.NewArgumentsError(`SemVer`, err.Error()))
				}
				return WrapSemVer(v)
			})
		},

		func(d eval.Dispatch) {
			d.Param(`SemVerHash`)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				hash := args[0].(*HashValue)
				major := hash.Get5(`major`, ZERO).(*IntegerValue).Int()
				minor := hash.Get5(`minor`, ZERO).(*IntegerValue).Int()
				patch := hash.Get5(`patch`, ZERO).(*IntegerValue).Int()
				preRelease := ``
				build := ``
				ev := hash.Get5(`prerelease`, nil)
				if ev != nil {
					preRelease = ev.String()
				}
				ev = hash.Get5(`build`, nil)
				if ev != nil {
					build = ev.String()
				}
				v, err := semver.NewVersion3(int(major), int(minor), int(patch), preRelease, build)
				if err != nil {
					panic(errors.NewArgumentsError(`SemVer`, err.Error()))
				}
				return WrapSemVer(v)
			})
		},
	)
}

func DefaultSemVerType() *SemVerType {
	return semVerType_DEFAULT
}

func NewSemVerType(vr semver.VersionRange) *SemVerType {
	if vr.Equals(semver.MatchAll) {
		return DefaultSemVerType()
	}
	return &SemVerType{vr}
}

func NewSemVerType2(limits ...eval.Value) *SemVerType {
	return NewSemVerType3(WrapValues(limits))
}

func NewSemVerType3(limits eval.List) *SemVerType {
	argc := limits.Len()
	if argc == 0 {
		return DefaultSemVerType()
	}

	if argc == 1 {
		if ranges, ok := limits.At(0).(eval.List); ok {
			return NewSemVerType3(ranges)
		}
	}

	var finalRange semver.VersionRange
	limits.EachWithIndex(func(arg eval.Value, idx int) {
		var rng semver.VersionRange
		str, ok := arg.(*StringValue)
		if ok {
			var err error
			rng, err = semver.ParseVersionRange(str.String())
			if err != nil {
				panic(errors.NewIllegalArgument(`SemVer[]`, idx, err.Error()))
			}
		} else {
			rv, ok := arg.(*SemVerRangeValue)
			if !ok {
				panic(NewIllegalArgumentType2(`SemVer[]`, idx, `Variant[String,SemVerRange]`, arg))
			}
			rng = rv.VersionRange()
		}
		if finalRange == nil {
			finalRange = rng
		} else {
			finalRange = finalRange.Merge(rng)
		}
	})
	return NewSemVerType(finalRange)
}

func (t *SemVerType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
}

func (t *SemVerType) Default() eval.Type {
	return semVerType_DEFAULT
}

func (t *SemVerType) Equals(o interface{}, g eval.Guard) bool {
	_, ok := o.(*SemVerType)
	return ok
}

func (t *SemVerType) Get(key string) (eval.Value, bool) {
	switch key {
	case `ranges`:
		return WrapValues(t.Parameters()), true
	default:
		return nil, false
	}
}

func (t *SemVerType) MetaType() eval.ObjectType {
	return SemVer_Type
}

func (t *SemVerType) Name() string {
	return `SemVer`
}

func (t *SemVerType) ReflectType(c eval.Context) (reflect.Type, bool) {
	return reflect.TypeOf(semver.Max), true
}

func (t *SemVerType)  CanSerializeAsString() bool {
  return true
}

func (t *SemVerType)  SerializationString() string {
	return t.String()
}


func (t *SemVerType) String() string {
	return eval.ToString2(t, NONE)
}

func (t *SemVerType) IsAssignable(o eval.Type, g eval.Guard) bool {
	if vt, ok := o.(*SemVerType); ok {
		return vt.vRange.IsAsRestrictiveAs(t.vRange)
	}
	return false
}

func (t *SemVerType) IsInstance(o eval.Value, g eval.Guard) bool {
	if v, ok := o.(*SemVerValue); ok {
		return t.vRange.Includes(v.Version())
	}
	return false
}

func (t *SemVerType) Parameters() []eval.Value {
	if t.vRange.Equals(semver.MatchAll) {
		return eval.EMPTY_VALUES
	}
	return []eval.Value{WrapString(t.vRange.String())}
}

func (t *SemVerType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *SemVerType) PType() eval.Type {
	return &TypeType{t}
}

func WrapSemVer(val semver.Version) *SemVerValue {
	return (*SemVerValue)(NewSemVerType(semver.ExactVersionRange(val)))
}

func (v *SemVerValue) Version() semver.Version {
	return v.vRange.StartVersion()
}

func (v *SemVerValue) Equals(o interface{}, g eval.Guard) bool {
	if ov, ok := o.(*SemVerValue); ok {
		return v.Version().Equals(ov.Version())
	}
	return false
}

func (v *SemVerValue) Reflect(c eval.Context) reflect.Value {
	return reflect.ValueOf(v.Version())
}

func (v *SemVerValue) ReflectTo(c eval.Context, dest reflect.Value) {
	rv := v.Reflect(c)
	if !rv.Type().AssignableTo(dest.Type()) {
		panic(eval.Error(eval.EVAL_ATTEMPT_TO_SET_WRONG_KIND, issue.H{`expected`: rv.Type().String(), `actual`: dest.Type().String()}))
	}
	dest.Set(rv)
}

func (v *SemVerValue)  CanSerializeAsString() bool {
  return true
}

func (v *SemVerValue)  SerializationString() string {
	return v.String()
}


func (v *SemVerValue) String() string {
	return v.Version().String()
}

func (v *SemVerValue) ToKey(b *bytes.Buffer) {
	b.WriteByte(1)
	b.WriteByte(HK_VERSION)
	v.Version().ToString(b)
}

func (v *SemVerValue) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	f := eval.GetFormat(s.FormatMap(), v.PType())
	val := v.Version().String()
	switch f.FormatChar() {
	case 's':
		f.ApplyStringFlags(b, val, f.IsAlt())
	case 'p':
		io.WriteString(b, `SemVer(`)
		utils.PuppetQuote(b, val)
		utils.WriteByte(b, ')')
	default:
		panic(s.UnsupportedFormat(v.PType(), `sp`, f))
	}
}

func (v *SemVerValue) PType() eval.Type {
	return (*SemVerType)(v)
}
