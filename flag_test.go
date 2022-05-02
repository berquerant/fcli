package fcli_test

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/berquerant/fcli"
	"github.com/stretchr/testify/assert"
)

type invalidCustomFlagType struct{}

func TestCustomFlagInvalid(t *testing.T) {
	_, ok := fcli.NewFlagFactory(reflect.TypeOf(&invalidCustomFlagType{}))
	assert.False(t, ok)
}

type customFlagTestcase struct {
	name        string
	sampleValue any // because allocates Flag by NewFlagFactory
	arg         string
	valueP      func(*testing.T, any)
}

func (s *customFlagTestcase) test(t *testing.T) {
	const (
		flagName    = "fname"
		flagSetName = "fset"
	)
	ff, ok := fcli.NewFlagFactory(reflect.TypeOf(s.sampleValue))
	if !assert.True(t, ok, "flag factory") {
		return
	}
	flg := ff(flagName)
	flgSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
	flg.AddFlag(flgSet)
	assert.Nil(t, flgSet.Parse([]string{
		fmt.Sprintf("-%s", flagName),
		s.arg,
	}), "parse")
	v, err := flg.Unwrap()
	if !assert.Nil(t, err) {
		return
	}
	s.valueP(t, v)
}

type customFlagStruct struct {
	value int
}

func (*customFlagStruct) UnmarshalFlag(_ string) (fcli.CustomFlagUnmarshaller, error) {
	return &customFlagStruct{
		value: 100,
	}, nil
}

type customFlagEnum int

func (customFlagEnum) UnmarshalFlag(_ string) (fcli.CustomFlagUnmarshaller, error) {
	return customFlagEnum(100), nil
}

type customFlagSlice []int

func (customFlagSlice) UnmarshalFlag(_ string) (fcli.CustomFlagUnmarshaller, error) {
	return customFlagSlice([]int{1, 2, 3}), nil
}

func TestCustomFlag(t *testing.T) {
	for _, tc := range []customFlagTestcase{
		{
			name:        "slice",
			sampleValue: customFlagSlice([]int{}),
			arg:         "arg",
			valueP: func(t *testing.T, v any) {
				x, ok := v.(customFlagSlice)
				if !assert.True(t, ok) {
					return
				}
				assert.Equal(t, customFlagSlice([]int{1, 2, 3}), x)
			},
		},
		{
			name:        "struct",
			sampleValue: &customFlagStruct{},
			arg:         "arg",
			valueP: func(t *testing.T, v any) {
				x, ok := v.(*customFlagStruct)
				if !assert.True(t, ok) {
					return
				}
				assert.Equal(t, 100, x.value)
			},
		},
		{
			name:        "enum",
			sampleValue: customFlagEnum(0),
			arg:         "arg",
			valueP: func(t *testing.T, v any) {
				x, ok := v.(customFlagEnum)
				if !assert.True(t, ok) {
					return
				}
				assert.Equal(t, customFlagEnum(100), x)
			},
		},
	} {
		tc := tc
		t.Run(tc.name, tc.test)
	}
}

type flagOutOfRangeTestcase struct {
	name string
	arg  string
}

type flagTestcase struct {
	name         string
	sampleValue  any    // because allocates Flag by NewFlagFactory
	arg          string // command-line arguments, like -`flagName` `arg`
	defaultValue any
	value        any
	outOfRange   []flagOutOfRangeTestcase
}

func (s *flagTestcase) test(t *testing.T) {
	const (
		flagName    = "fname"
		flagSetName = "fset"
	)

	arg := func(a string) []string {
		f := fmt.Sprintf("-%s", flagName)
		if _, ok := s.sampleValue.(bool); ok {
			// bool flag specifies like `-boolFlag`
			return []string{f}
		}
		return []string{f, a}
	}

	ff, ok := fcli.NewFlagFactory(reflect.TypeOf(s.sampleValue))
	if !assert.True(t, ok, "flag factory") {
		return
	}
	flg := ff(flagName)
	flgSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
	flg.AddFlag(flgSet)

	t.Run("default value", func(t *testing.T) {
		assert.Nil(t, flgSet.Parse([]string{}), "parse")
		v, err := flg.Unwrap()
		assert.Nil(t, err)
		assert.Equal(t, s.defaultValue, v)
	})

	t.Run("arg", func(t *testing.T) {
		assert.Nil(t, flgSet.Parse(arg(s.arg)), "parse")
		v, err := flg.Unwrap()
		assert.Nil(t, err)
		assert.Equal(t, s.value, v)
	})

	for _, tc := range s.outOfRange {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Nil(t, flgSet.Parse(arg(tc.arg)), "parse")
			_, err := flg.Unwrap()
			assert.ErrorIs(t, err, fcli.ErrValueOutOfRange)
		})
	}
}

func TestFlag(t *testing.T) {
	for _, tc := range []flagTestcase{
		{
			name:         "bool",
			sampleValue:  false,
			defaultValue: false,
			value:        true,
		},
		{
			name:         "int",
			sampleValue:  int(0),
			arg:          "2",
			defaultValue: 0,
			value:        2,
		},
		{
			name:         "int8",
			sampleValue:  int8(0),
			arg:          "5",
			defaultValue: int8(0),
			value:        int8(5),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  "128",
				},
				{
					name: "too small",
					arg:  "-129",
				},
			},
		},
		{
			name:         "int16",
			sampleValue:  int16(0),
			arg:          "129",
			defaultValue: int16(0),
			value:        int16(129),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  "32768",
				},
				{
					name: "too small",
					arg:  "-32769",
				},
			},
		},
		{
			name:         "int32",
			sampleValue:  int32(0),
			arg:          "32768",
			defaultValue: int32(0),
			value:        int32(32768),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  "2147483648",
				},
				{
					name: "too small",
					arg:  "-2147483649",
				},
			},
		},
		{
			name:         "int64",
			sampleValue:  int64(0),
			arg:          "2147483648",
			defaultValue: int64(0),
			value:        int64(2147483648),
		},
		{
			name:         "uint",
			sampleValue:  uint(0),
			arg:          "11",
			defaultValue: uint(0),
			value:        uint(11),
		},
		{
			name:         "uint8",
			sampleValue:  uint8(0),
			arg:          "11",
			defaultValue: uint8(0),
			value:        uint8(11),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  "256",
				},
			},
		},
		{
			name:         "uint16",
			sampleValue:  uint16(0),
			arg:          "256",
			defaultValue: uint16(0),
			value:        uint16(256),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  "65536",
				},
			},
		},
		{
			name:         "uint32",
			sampleValue:  uint32(0),
			arg:          "65536",
			defaultValue: uint32(0),
			value:        uint32(65536),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  "4294967296",
				},
			},
		},
		{
			name:         "uint64",
			sampleValue:  uint64(0),
			arg:          "4294967296",
			defaultValue: uint64(0),
			value:        uint64(4294967296),
		},
		{
			name:         "float32",
			sampleValue:  float32(0),
			arg:          "1.2",
			defaultValue: float32(0),
			value:        float32(1.2),
			outOfRange: []flagOutOfRangeTestcase{
				{
					name: "too big",
					arg:  fmt.Sprint(4e+38),
				},
				{
					name: "too small",
					arg:  fmt.Sprint(-4e+38),
				},
			},
		},
		{
			name:         "float64",
			sampleValue:  float64(0),
			arg:          fmt.Sprint(4e+38),
			defaultValue: float64(0),
			value:        4e+38,
		},
		{
			name:         "string",
			sampleValue:  "",
			arg:          "str",
			defaultValue: "",
			value:        "str",
		},
	} {
		tc := tc
		t.Run(tc.name, tc.test)
	}
}
