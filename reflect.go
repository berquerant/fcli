package fcli

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"runtime"
	"strings"
)

var (
	ErrNotFunction       = errors.New("not function")
	ErrCannotCutFuncDecl = errors.New("cannot cut func decl")
	ErrInvalidFuncInfo   = errors.New("invalid func info")
)

// FuncName represents the function name and the location.
type FuncName struct {
	fullname string
	file     string
	line     int
}

func (s *FuncName) FullName() string { return s.fullname }
func (s *FuncName) File() string     { return s.file }
func (s *FuncName) Line() int        { return s.line }
func (s *FuncName) String() string {
	ss := strings.Split(s.fullname, ".")
	return ss[len(ss)-1]
}

// GetFuncName returns the function name.
// Returns ErrNotFunction if f is not a function.
// Note: this is not for method, literal.
func GetFuncName(f any) (*FuncName, error) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return nil, ErrNotFunction
	}
	ptr := reflect.ValueOf(f).Pointer()
	fu := runtime.FuncForPC(ptr)
	file, line := fu.FileLine(ptr)
	return &FuncName{
		fullname: fu.Name(),
		file:     file,
		line:     line,
	}, nil
}

// FileLines is the set of the lines of the file.
type FileLines interface {
	Filename() string
	Line(lineNumber int) (string, bool)
}

type fileLines struct {
	filename string
	lines    []string
}

func (s *fileLines) Filename() string { return s.filename }
func (s *fileLines) Line(lineNumber int) (string, bool) {
	idx := lineNumber - 1
	if idx < 0 || idx >= len(s.lines) {
		return "", false
	}
	return s.lines[idx], true
}

func NewFileLines(file string) (FileLines, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines := []string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &fileLines{
		filename: file,
		lines:    lines,
	}, nil
}

type FuncDeclCutter interface {
	// CutFuncDecl cuts out a func decl.
	// Include doc but comment marker // only.
	CutFuncDecl() (string, error)
}

func NewFuncDeclCutter(file FileLines, funcHeadLineNumber int) FuncDeclCutter {
	return &funcDeclCutter{
		file: file,
		line: funcHeadLineNumber,
	}
}

type funcDeclCutter struct {
	file FileLines
	line int
}

func (s *funcDeclCutter) CutFuncDecl() (string, error) {
	wrapErr := NewErrorWrapperBuilder().
		Err(ErrCannotCutFuncDecl).
		Msg("%s line %d", s.file.Filename(), s.line).
		Build()

	if _, ok := s.findFuncHead(); !ok {
		return "", wrapErr("no func decl head")
	}

	body, ok := s.findFuncBody()
	if !ok {
		return "", wrapErr("cannot find body")
	}
	doc := s.findFuncDocComment()

	return strings.Join(append(doc, body...), "\n"), nil
}

func (s *funcDeclCutter) findFuncHead() (string, bool) {
	pivot, ok := s.file.Line(s.line)
	if !ok || !strings.HasPrefix(pivot, "func ") {
		return "", false
	}
	return pivot, true
}

func (s *funcDeclCutter) findFuncBody() ([]string, bool) {
	var (
		leftBraces  int
		rightBraces int
		idx         = s.line
	)

	for {
		line, ok := s.file.Line(idx)
		if !ok {
			return nil, false
		}
		leftBraces += strings.Count(line, "{")
		rightBraces += strings.Count(line, "}")
		if leftBraces > 0 && rightBraces > 0 && leftBraces == rightBraces {
			break
		}
		idx++
	}

	var (
		j     int
		lines = make([]string, idx-s.line+1)
	)
	for i := s.line; i <= idx; i++ {
		lines[j], _ = s.file.Line(i)
		j++
	}
	return lines, true
}

func (s *funcDeclCutter) findFuncDocComment() []string {
	startIdx := s.line
	for {
		x, ok := s.file.Line(startIdx - 1)
		if !ok || !strings.HasPrefix(x, "//") {
			break
		}
		startIdx--
	}

	var (
		j        int
		comments = make([]string, s.line-startIdx)
	)
	for i := startIdx; i < s.line; i++ {
		comments[j], _ = s.file.Line(i)
		j++
	}
	return comments
}

// FuncParam is an input parameter of the function.
type FuncParam interface {
	// Name returns the name of the parameter.
	Name() string
}

type FuncInfo interface {
	// Name returns the name of the function.
	Name() string
	// Doc returns the comments of the function.
	Doc() string
	NumIn() int
	In(int) FuncParam
}

type funcParam struct {
	name string
}

func (s *funcParam) Name() string { return s.name }

// NewFuncInfo parses src and generate FuncInfo.
// src is func decl, like:
//
//   func targetFunc() int { return 0 }
//
// Returns ErrInvalidFuncInfo if parse failed.
func NewFuncInfo(src string) (FuncInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fmt.Sprintf("package tmp\n%s", src), parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("%w %v %s", ErrInvalidFuncInfo, err, src)
	}

	var (
		found bool
		fdecl *ast.FuncDecl
	)
	for _, decl := range f.Decls {
		if found {
			break
		}
		ast.Inspect(decl, func(node ast.Node) bool {
			if x, ok := node.(*ast.FuncDecl); ok {
				fdecl = x
				found = true
			}
			return false
		})
	}

	if !found {
		return nil, fmt.Errorf("%w func decl not found", ErrInvalidFuncInfo)
	}
	inNames := []string{}
	for _, names := range fdecl.Type.Params.List {
		for _, name := range names.Names {
			inNames = append(inNames, name.Name)
		}
	}
	return &funcInfo{
		decl:    fdecl,
		inNames: inNames,
	}, nil
}

type funcInfo struct {
	decl    *ast.FuncDecl
	inNames []string
}

func (s *funcInfo) Doc() string  { return s.decl.Doc.Text() }
func (s *funcInfo) Name() string { return s.decl.Name.Name }
func (s *funcInfo) NumIn() int   { return len(s.inNames) }
func (s *funcInfo) In(i int) FuncParam {
	if i < 0 || i >= len(s.inNames) {
		return nil
	}
	return &funcParam{
		name: s.inNames[i],
	}
}
