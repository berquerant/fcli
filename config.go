package fcli

import "flag"

type Config struct {
	ErrorHandling flag.ErrorHandling
}

func (s *Config) Apply(opt ...Option) {
	for _, x := range opt {
		x(s)
	}
}

type Option func(*Config)

func WithErrorHandling(h flag.ErrorHandling) Option {
	return func(config *Config) {
		config.ErrorHandling = h
	}
}
