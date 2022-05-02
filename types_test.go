package fcli_test

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/berquerant/fcli"
	"github.com/stretchr/testify/assert"
)

func testTargetFunctionCallTarget() {}

func TestTargetFunctionCallName(t *testing.T) {
	s, err := fcli.NewTargetFunction(testTargetFunctionCallTarget)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, s.Name(), "testTargetFunctionCallTarget")
}

type targetFunctionTestcaseResult struct {
	result []any
	sync.Mutex
}

var targetFunctionTestcaseResultInstance targetFunctionTestcaseResult

func getTargetFunctionTestcaseResult() []any  { return targetFunctionTestcaseResultInstance.result }
func setTargetFunctionTestcaseResult(v []any) { targetFunctionTestcaseResultInstance.result = v }

type targetFunctionCallTestcase struct {
	name      string
	f         any // func(...)
	args      []string
	wantArgsP func(*testing.T, []any)
	newErr    error
	callErr   error
}

func variadic(_ int, _ ...any)            {}
func withOutputParam() error              { return nil }
func unsupportedInputParamType(_ uintptr) {}
func withoutInputParams()                 {}
func singleIntInput(i int) {
	setTargetFunctionTestcaseResult([]any{i})
}
func allBasic(b bool,
	i int, i8 int8, i16 int16, i32 int32, i64 int64,
	u uint, u8 uint8, u16 uint16, u32 uint32, u64 uint64,
	f32 float32, f64 float64,
	str string) {
	setTargetFunctionTestcaseResult([]any{
		b,
		i, i8, i16, i32, i64,
		u, u8, u16, u32, u64,
		f32, f64,
		str,
	})
}

type stringList struct {
	list []string
}

func (*stringList) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	return &stringList{
		list: strings.Split(v, ","),
	}, nil
}

type failUnmarshaller struct{}

func (*failUnmarshaller) UnmarshalFlag(_ string) (fcli.CustomFlagUnmarshaller, error) {
	return nil, fmt.Errorf("fail unmarshaller")
}

func customStringList(list *stringList) {
	setTargetFunctionTestcaseResult([]any{list})
}

func customFlagFailure(v *failUnmarshaller) {}

func int8LimitCheck(i8 int8) {}

func TestTargetFunctionCall(t *testing.T) {
	for _, tc := range []targetFunctionCallTestcase{
		{
			name:   "not a function",
			f:      1,
			newErr: fcli.ErrBadTargetFunction,
		},
		{
			name:   "variadic",
			f:      variadic,
			newErr: fcli.ErrBadTargetFunction,
		},
		{
			name:   "with output params",
			f:      withOutputParam,
			newErr: fcli.ErrBadTargetFunction,
		},
		{
			name:   "unsupported input param type",
			f:      unsupportedInputParamType,
			newErr: fcli.ErrBadTargetFunction,
		},
		{
			name:      "no input params",
			f:         withoutInputParams,
			wantArgsP: func(_ *testing.T, _ []any) {},
		},
		{
			name: "int",
			f:    singleIntInput,
			args: []string{
				"-i", "1",
			},
			wantArgsP: func(t *testing.T, v []any) {
				assert.Equal(t, []any{1}, v)
			},
		},
		{
			name: "int default",
			f:    singleIntInput,
			args: []string{},
			wantArgsP: func(t *testing.T, v []any) {
				assert.Equal(t, []any{0}, v)
			},
		},
		{
			name: "int call failure",
			f:    singleIntInput,
			args: []string{
				"-i", "INT",
			},
			callErr: fcli.ErrCallFailure,
		},
		{
			name: "all basic",
			f:    allBasic,
			args: []string{
				"-b",
				"-i", "1",
				"-i8", "8",
				"-i16", "16",
				"-i32", "32",
				"-i64", "64",
				"-u", "2",
				"-u8", "9",
				"-u16", "17",
				"-u32", "33",
				"-u64", "65",
				"-f32", "37.1",
				"-f64", "67.2",
				"-str", "all-basic",
			},
			wantArgsP: func(t *testing.T, v []any) {
				want := []any{
					true,
					int(1),
					int8(8),
					int16(16),
					int32(32),
					int64(64),
					uint(2),
					uint8(9),
					uint16(17),
					uint32(33),
					uint64(65),
					float32(37.1),
					float64(67.2),
					"all-basic",
				}
				assert.Equal(t, want, v)
			},
		},
		{
			name: "string list",
			f:    customStringList,
			args: []string{
				"-list", "a,b,c",
			},
			wantArgsP: func(t *testing.T, v []any) {
				if !assert.Equal(t, 1, len(v), v) {
					return
				}
				x := v[0]
				y, ok := x.(*stringList)
				if !assert.True(t, ok, x) {
					return
				}
				assert.Equal(t, []string{"a", "b", "c"}, y.list)
			},
		},
		{
			name:    "custom flag failure",
			f:       customFlagFailure,
			args:    []string{"-v", "fail"},
			callErr: fcli.ErrCallFailure,
		},
		{
			name: "int8 limit",
			f:    int8LimitCheck,
			args: []string{
				"-i8", "128",
			},
			callErr: fcli.ErrCallFailure,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			targetFunctionTestcaseResultInstance.Lock()
			defer targetFunctionTestcaseResultInstance.Unlock()

			s, err := fcli.NewTargetFunction(tc.f, fcli.WithErrorHandling(flag.ContinueOnError))
			assert.ErrorIs(t, err, tc.newErr, "new error")
			if tc.newErr != nil {
				t.Logf("new error %v", err)
				return
			}
			err = s.Call(tc.args)
			assert.ErrorIs(t, err, tc.callErr, "call error")
			if tc.callErr != nil {
				t.Logf("call error %v", err)
				return
			}
			tc.wantArgsP(t, getTargetFunctionTestcaseResult())
			setTargetFunctionTestcaseResult(nil)
		})
	}
}
