package fcli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/berquerant/fcli"
	"github.com/stretchr/testify/assert"
)

func TestGetFuncName(t *testing.T) {
	for _, tc := range []struct {
		name string
		f    any
		want string
		err  error
	}{
		{
			name: "not function",
			f:    1,
			err:  fcli.ErrNotFunction,
		},
		{
			name: "func",
			f:    fcli.GetFuncName,
			want: "GetFuncName",
		},
		{
			name: "assert",
			f:    assert.Equal,
			want: "Equal",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := fcli.GetFuncName(tc.f)
			assert.Equal(t, tc.err, err)
			if tc.err != nil {
				return
			}
			assert.Equal(t, tc.want, got.String())
		})
	}
}

type mockFileLines struct {
	lines []string
}

func (*mockFileLines) Filename() string { return "" }
func (s *mockFileLines) Line(lineNumber int) (string, bool) {
	i := lineNumber - 1
	if i < 0 || i >= len(s.lines) {
		return "", false
	}
	return s.lines[i], true
}

func TestFuncDeclCutter(t *testing.T) {
	for _, tc := range []struct {
		name string
		file string
		line int
		want string
		err  error
	}{
		{
			name: "not function",
			file: `package tmp
const = 1`,
			line: 2,
			err:  fcli.ErrCannotCutFuncDecl,
		},
		{
			name: "incomplete function",
			file: `package tmp
func f() {
`,
			line: 2,
			err:  fcli.ErrCannotCutFuncDecl,
		},
		{
			name: "inline",
			file: `package tmp
func f() {}`,
			line: 2,
			want: `func f() {}`,
		},
		{
			name: "func",
			file: `package tmp
func f(name string) {
  println(name)
}
func g() {}`,
			line: 2,
			want: `func f(name string) {
  println(name)
}`,
		},
		{
			name: "inline with doc",
			file: `package tmp
// A
// B
func withDoc() {}
`,
			line: 4,
			want: `// A
// B
func withDoc() {}`,
		},
		{
			name: "with doc",
			file: `package tmp
// A
//
// B
func withDoc(name string) int {
  println(name)
  return 0
}
`,
			line: 5,
			want: `// A
//
// B
func withDoc(name string) int {
  println(name)
  return 0
}`,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			lines := &mockFileLines{
				lines: strings.Split(tc.file, "\n"),
			}
			s := fcli.NewFuncDeclCutter(lines, tc.line)
			got, err := s.CutFuncDecl()
			assert.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				t.Logf("got error %v", err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

type funcInfoTestcase struct {
	name        string
	src         string
	funcName    string
	doc         string
	wantInNames []string
	err         error
}

func (s *funcInfoTestcase) test(t *testing.T) {
	got, err := fcli.NewFuncInfo(s.src)
	assert.ErrorIs(t, err, s.err)
	if s.err != nil {
		t.Logf("got error %v", err)
		return
	}
	assert.Equal(t, s.funcName, got.Name())
	assert.Equal(t, s.doc, got.Doc())
	if !assert.Equal(t, len(s.wantInNames), got.NumIn()) {
		return
	}
	for i := 0; i < got.NumIn(); i++ {
		assert.Equal(t, s.wantInNames[i], got.In(i).Name(), fmt.Sprintf("name[%d]", i))
	}
}

func TestFuncInfo(t *testing.T) {
	for _, tc := range []funcInfoTestcase{
		{
			name: "no func",
			src:  ``,
			err:  fcli.ErrInvalidFuncInfo,
		},
		{
			name:        "inline subroutine",
			src:         `func subroutine() {}`,
			funcName:    "subroutine",
			wantInNames: []string{},
		},
		{
			name:        "inline func",
			src:         `func add(x, y int) int { return x + y }`,
			funcName:    "add",
			wantInNames: []string{"x", "y"},
		},
		{
			name: "func",
			src: `func query(category string, x, y, z int) string {
  if category == "sum" {
    return fmt.Sprint(x + y + z)
  }
  return ""
}`,
			funcName:    "query",
			wantInNames: []string{"category", "x", "y", "z"},
		},
		{
			name: "inline with doc",
			src: `// Hello!
func Hello() {}`,
			funcName:    "Hello",
			wantInNames: []string{},
			doc: `Hello!
`,
		},
		{
			name: "doc",
			src: `// Multiline
// comment
func Multiline() {
  println("Hello")
}`,
			funcName:    "Multiline",
			wantInNames: []string{},
			doc: `Multiline
comment
`,
		},
	} {
		tc := tc
		t.Run(tc.name, tc.test)
	}
}
