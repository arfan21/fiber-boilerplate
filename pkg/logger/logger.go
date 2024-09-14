package logger

import (
	"context"
	"os"
	"sync"

	"github.com/arfan21/fiber-boilerplate/config"
	"github.com/arfan21/otelzerolog"
	"github.com/rs/zerolog"
)

var loggerInstance zerolog.Logger
var once sync.Once

func Log(ctx context.Context) *zerolog.Logger {
	once.Do(func() {
		multi := zerolog.MultiLevelWriter(os.Stdout)
		loggerInstance = zerolog.New(multi).With().Timestamp().Logger()

		if config.Get().Env == "dev" {
			loggerInstance = loggerInstance.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		if config.Get().Otel.Enabled {
			loggerInstance = loggerInstance.Hook(otelzerolog.NewHook(config.Get().Service.Name))
		}

	})

	newlogger := loggerInstance.With().Ctx(ctx).Logger()
	return &newlogger
}
