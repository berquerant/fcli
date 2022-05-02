package fcli

import (
	"flag"
	"fmt"
	"math"
	"reflect"
)

// Flag is a command-line flag.
type Flag interface {
	// Name is the flag name.
	Name() string
	// AddFlag defines the flag in the flag set.
	AddFlag(flagSet *flag.FlagSet)
	// Unwrap returns the flag value.
	// Default value is the zero value.
	Unwrap() (any, error)
	// ReflectValue returns the flag value for reflection.
	// Default value is the zero value.
	ReflectValue() (reflect.Value, error)
}

type baseFlag struct {
	name string
}

func (s *baseFlag) Name() string { return s.name }

func newBaseFlag(name string) *baseFlag {
	return &baseFlag{
		name: name,
	}
}

// IntFlag is a flag for int.
type IntFlag struct {
	*baseFlag
	value *int
}

func NewIntFlag(name string) Flag {
	return &IntFlag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *IntFlag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.Int(s.name, 0, "") }
func (s *IntFlag) Value() (int, error)           { return *s.value, nil }
func (s *IntFlag) Unwrap() (any, error)          { return s.Value() }

func (s *IntFlag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

var (
	// ErrValueOutOfRange is the error returned if the parsed value from the command-line arguments
	// is the out of the range.
	ErrValueOutOfRange = fmt.Errorf("value out of range")
)

type sizedIntFlag struct {
	*IntFlag
	min, max int
}

func newSizedIntFlag(name string, min, max int) *sizedIntFlag {
	return &sizedIntFlag{
		IntFlag: &IntFlag{
			baseFlag: newBaseFlag(name),
		},
		min: min,
		max: max,
	}
}

func (s *sizedIntFlag) Value() (int, error) {
	v := *s.value
	if v < s.min || v > s.max {
		return 0, ErrValueOutOfRange
	}
	return v, nil
}

func (s *sizedIntFlag) Unwrap() (any, error) { return s.Value() }

// Int8Flag is the flag for int8.
type Int8Flag struct {
	*sizedIntFlag
}

func NewInt8Flag(name string) Flag {
	return &Int8Flag{
		sizedIntFlag: newSizedIntFlag(name, math.MinInt8, math.MaxInt8),
	}
}

func (s *Int8Flag) Unwrap() (any, error) { return s.Value() }

func (s *Int8Flag) Value() (int8, error) {
	v, err := s.sizedIntFlag.Value()
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

func (s *Int8Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Int16Flag is the flag for int16.
type Int16Flag struct {
	*sizedIntFlag
}

func NewInt16Flag(name string) Flag {
	return &Int16Flag{
		sizedIntFlag: newSizedIntFlag(name, math.MinInt16, math.MaxInt16),
	}
}

func (s *Int16Flag) Unwrap() (any, error) { return s.Value() }

func (s *Int16Flag) Value() (int16, error) {
	v, err := s.sizedIntFlag.Value()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (s *Int16Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Int32Flag is the flag for int32.
type Int32Flag struct {
	*sizedIntFlag
}

func NewInt32Flag(name string) Flag {
	return &Int32Flag{
		sizedIntFlag: newSizedIntFlag(name, math.MinInt32, math.MaxInt32),
	}
}

func (s *Int32Flag) Unwrap() (any, error) { return s.Value() }

func (s *Int32Flag) Value() (int32, error) {
	v, err := s.sizedIntFlag.Value()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (s *Int32Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Int64Flag is the flag for int64.
type Int64Flag struct {
	*baseFlag
	value *int64
}

func NewInt64Flag(name string) Flag {
	return &Int64Flag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *Int64Flag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.Int64(s.name, 0, "") }
func (s *Int64Flag) Value() (int64, error)         { return *s.value, nil }
func (s *Int64Flag) Unwrap() (any, error)          { return s.Value() }
func (s *Int64Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// UintFlag is the flag for uint.
type UintFlag struct {
	*baseFlag
	value *uint
}

func NewUintFlag(name string) Flag {
	return &UintFlag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *UintFlag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.Uint(s.name, 0, "") }
func (s *UintFlag) Value() (uint, error)          { return *s.value, nil }
func (s *UintFlag) Unwrap() (any, error)          { return s.Value() }
func (s *UintFlag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

type sizedUintFlag struct {
	*UintFlag
	max uint
}

func newSizedUintFlag(name string, max uint) *sizedUintFlag {
	return &sizedUintFlag{
		UintFlag: &UintFlag{
			baseFlag: newBaseFlag(name),
		},
		max: max,
	}
}

func (s *sizedUintFlag) Value() (uint, error) {
	v := *s.value
	if v > s.max {
		return 0, ErrValueOutOfRange
	}
	return v, nil
}

func (s *sizedUintFlag) Unwrap() (any, error) { return s.Value() }

// Uint8Flag is the flag for uint8.
type Uint8Flag struct {
	*sizedUintFlag
}

func NewUint8Flag(name string) Flag {
	return &Uint8Flag{
		sizedUintFlag: newSizedUintFlag(name, math.MaxUint8),
	}
}

func (s *Uint8Flag) Unwrap() (any, error) { return s.Value() }

func (s *Uint8Flag) Value() (uint8, error) {
	v, err := s.sizedUintFlag.Value()
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

func (s *Uint8Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Uint16Flag is the flag for uint16.
type Uint16Flag struct {
	*sizedUintFlag
}

func NewUint16Flag(name string) Flag {
	return &Uint16Flag{
		sizedUintFlag: newSizedUintFlag(name, math.MaxUint16),
	}
}

func (s *Uint16Flag) Unwrap() (any, error) { return s.Value() }

func (s *Uint16Flag) Value() (uint16, error) {
	v, err := s.sizedUintFlag.Value()
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

func (s *Uint16Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Uint32Flag is the flag for uint32.
type Uint32Flag struct {
	*sizedUintFlag
}

func NewUint32Flag(name string) Flag {
	return &Uint32Flag{
		sizedUintFlag: newSizedUintFlag(name, math.MaxUint32),
	}
}

func (s *Uint32Flag) Unwrap() (any, error) { return s.Value() }

func (s *Uint32Flag) Value() (uint32, error) {
	v, err := s.sizedUintFlag.Value()
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func (s *Uint32Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Uint64Flag is the flag for uint64.
type Uint64Flag struct {
	*baseFlag
	value *uint64
}

func NewUint64Flag(name string) Flag {
	return &Uint64Flag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *Uint64Flag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.Uint64(s.name, 0, "") }
func (s *Uint64Flag) Value() (uint64, error)        { return *s.value, nil }
func (s *Uint64Flag) Unwrap() (any, error)          { return s.Value() }
func (s *Uint64Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// BoolFlag is the flag for bool.
type BoolFlag struct {
	*baseFlag
	value *bool
}

func NewBoolFlag(name string) Flag {
	return &BoolFlag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *BoolFlag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.Bool(s.name, false, "") }
func (s *BoolFlag) Value() (bool, error)          { return *s.value, nil }
func (s *BoolFlag) Unwrap() (any, error)          { return s.Value() }
func (s *BoolFlag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Float64Flag is the flag for float64.
type Float64Flag struct {
	*baseFlag
	value *float64
}

func NewFloat64Flag(name string) Flag {
	return &Float64Flag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *Float64Flag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.Float64(s.name, 0, "") }
func (s *Float64Flag) Value() (float64, error)       { return *s.value, nil }
func (s *Float64Flag) Unwrap() (any, error)          { return s.Value() }
func (s *Float64Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

// Float32Flag is the flag for float32.
type Float32Flag struct {
	*Float64Flag
}

func NewFloat32Flag(name string) Flag {
	return &Float32Flag{
		Float64Flag: &Float64Flag{
			baseFlag: newBaseFlag(name),
		},
	}
}

func (s *Float32Flag) Value() (float32, error) {
	v, err := s.Float64Flag.Value()
	if err != nil {
		return 0, err
	}
	if v < -math.MaxFloat32 || v > math.MaxFloat32 {
		return 0, ErrValueOutOfRange
	}
	return float32(v), nil
}

func (s *Float32Flag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

func (s *Float32Flag) Unwrap() (any, error) { return s.Value() }

// StringFlag is the flag for string.
type StringFlag struct {
	*baseFlag
	value *string
}

func NewStringFlag(name string) Flag {
	return &StringFlag{
		baseFlag: newBaseFlag(name),
	}
}

func (s *StringFlag) AddFlag(flagSet *flag.FlagSet) { s.value = flagSet.String(s.name, "", "") }
func (s *StringFlag) Value() (string, error)        { return *s.value, nil }
func (s *StringFlag) Unwrap() (any, error)          { return s.Value() }
func (s *StringFlag) ReflectValue() (reflect.Value, error) {
	v, err := s.Value()
	return reflect.ValueOf(v), err
}

var (
	// ErrInvalidCustomFlag is the error returned if failed to parse the value of the custom flag.
	ErrCannotUnmarshalCustomFlag = fmt.Errorf("cannot unmarshal custom flag")
	// ErrInvalidCustomFlag is the error returned if the type of the flag is not proper.
	ErrInvalidCustomFlag = fmt.Errorf("invalid custom flag")
)

// CustomFlagUnmarshaller should be implemented by the type for CustomFlag.
type CustomFlagUnmarshaller interface {
	// UnmarshalFlag converts string into the value.
	// The value must be the implementing type.
	UnmarshalFlag(string) (CustomFlagUnmarshaller, error)
}

// CustomFlagZeroer provides zero value.
type CustomFlagZeroer interface {
	FlagZero() CustomFlagUnmarshaller
}

// CustomFlag is the flag for types that implement CustomFlagUnmarshaller.
type CustomFlag struct {
	*baseFlag
	typ   reflect.Type
	value any
}

// NewCustomFlag returns the new CustomFlag.
// typ is the type that implements CustomFlagUnmarshaller.
// struct should be passed as pointer.
// If typ implements CustomFlagZeroer, use it as the default value of the flag.
func NewCustomFlag(name string, typ reflect.Type) (Flag, error) {
	v := reflect.Zero(typ).Interface()
	if _, ok := v.(CustomFlagUnmarshaller); ok {
		f := &CustomFlag{
			baseFlag: newBaseFlag(name),
			typ:      typ,
		}
		if x, ok := v.(CustomFlagZeroer); ok {
			f.value = x.FlagZero()
		}
		return f, nil
	}
	return nil, fmt.Errorf("%w %s", ErrInvalidCustomFlag, name)
}

func (s *CustomFlag) AddFlag(flagSet *flag.FlagSet) {
	flagSet.Func(s.name, "", s.parse)
}
func (s *CustomFlag) parse(v string) error {
	z := reflect.Zero(s.typ).Interface()
	m, ok := z.(CustomFlagUnmarshaller)
	if !ok {
		return fmt.Errorf("%w type assertion %#v %s %s", ErrCannotUnmarshalCustomFlag, z, s.name, v)
	}
	p, err := m.UnmarshalFlag(v)
	if err != nil {
		return fmt.Errorf("%w UnmarshalFlag() %s %s %v", ErrCannotUnmarshalCustomFlag, s.name, v, err)
	}
	s.value = p
	return nil
}
func (s *CustomFlag) Unwrap() (any, error)                 { return s.value, nil }
func (s *CustomFlag) ReflectValue() (reflect.Value, error) { return reflect.ValueOf(s.value), nil }

type FlagFactory func(name string) Flag

// NewFlagFactory returns the proper flag for the v.
func NewFlagFactory(v any) (FlagFactory, bool) {
	newCustomFlag := func(typ reflect.Type) (FlagFactory, bool) {
		if _, err := NewCustomFlag("", typ); err != nil {
			return nil, false
		}
		return func(name string) Flag {
			f, _ := NewCustomFlag(name, typ)
			return f
		}, true
	}

	t := func() reflect.Type {
		if rt, ok := v.(reflect.Type); ok {
			return rt
		}
		return reflect.TypeOf(v)
	}()

	if f, ok := newCustomFlag(t); ok {
		return f, true
	}

	switch t.Kind() {
	case reflect.Bool:
		return NewBoolFlag, true
	case reflect.Int:
		return NewIntFlag, true
	case reflect.Int8:
		return NewInt8Flag, true
	case reflect.Int16:
		return NewInt16Flag, true
	case reflect.Int32:
		return NewInt32Flag, true
	case reflect.Int64:
		return NewInt64Flag, true
	case reflect.Uint:
		return NewUintFlag, true
	case reflect.Uint8:
		return NewUint8Flag, true
	case reflect.Uint16:
		return NewUint16Flag, true
	case reflect.Uint32:
		return NewUint32Flag, true
	case reflect.Uint64:
		return NewUint64Flag, true
	case reflect.Float32:
		return NewFloat32Flag, true
	case reflect.Float64:
		return NewFloat64Flag, true
	case reflect.String:
		return NewStringFlag, true
	default:
		return nil, false
	}
}
