package logger

import (
	"github.com/yakumioto/otelslog"
	"log/slog"
	"os"
	"runtime/debug"
)

type Config struct {
	level              slog.Level
	IncludeProgramInfo bool
}

func MustLogger(config Config) (*slog.Logger, error) {
	log := slog.New(
		otelslog.NewHandler(slog.NewJSONHandler(os.Stdout, nil)),
	)

	if config.IncludeProgramInfo {
		buildInfo, _ := debug.ReadBuildInfo()

		log = log.With(slog.Group("program_info", "os", os.Getpid(), "go_version", buildInfo.GoVersion))
	}

	slog.SetDefault(log)

	return log, nil
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
