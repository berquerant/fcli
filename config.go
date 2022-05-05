package fcli

import (
	"github.com/berquerant/fcli/internal/logger"
)

//go:generate go run github.com/berquerant/goconfig@latest -type "flag.ErrorHandling,CommandName|string" -option -output config_generated.go -configOption Option

func SetVerboseLevel(level int) {
	switch {
	case level <= 0:
		logger.SetLevel(logger.Lsilent)
	case level == 1:
		logger.SetLevel(logger.Linfo)
	case level == 2:
		logger.SetLevel(logger.Ldebug)
	case level >= 3:
		logger.SetLevel(logger.Ltrace)
	}
}
