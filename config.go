package fcli

import "flag"

type ConfigItem[T any] interface {
	// Get returns the stored value.
	// Returns the default value if Set is not called.
	Get() T
	// Set stores the value.
	Set(T)
}

type configItemImpl[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func NewConfigItem[T any](defaultValue T) ConfigItem[T] {
	return &configItemImpl[T]{
		defaultValue: defaultValue,
	}
}

func (s *configItemImpl[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *configItemImpl[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}

type ConfigBuilder struct {
	errorHandling flag.ErrorHandling
	commandName   string
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (s *ConfigBuilder) ErrorHandling(h flag.ErrorHandling) *ConfigBuilder {
	s.errorHandling = h
	return s
}
func (s *ConfigBuilder) CommandName(name string) *ConfigBuilder {
	s.commandName = name
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		ErrorHandling: NewConfigItem(s.errorHandling),
		CommandName:   NewConfigItem(s.commandName),
	}
}

type Config struct {
	ErrorHandling ConfigItem[flag.ErrorHandling]
	CommandName   ConfigItem[string]
}

func (s *Config) Apply(opt ...Option) {
	for _, x := range opt {
		x(s)
	}
}

type Option func(*Config)

// WithErrorHandling changes flag error handling.
// Default is flag.ExitOnError.
func WithErrorHandling(h flag.ErrorHandling) Option {
	return func(config *Config) {
		config.ErrorHandling.Set(h)
	}
}

// WithCommandName changes the command name.
// Default is the given function name.
func WithCommandName(name string) Option {
	return func(config *Config) {
		config.CommandName.Set(name)
	}
}
