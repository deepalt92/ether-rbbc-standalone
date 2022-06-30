package tracer

import (
	"ether-rbbc/log"
	"go.uber.org/zap"
)

type Tracer struct {
	shouldTrace bool
	Logger      zap.SugaredLogger
}

func NewTracer(cfg Config) (Tracer, error) {
	logger, err := log.New(cfg.LogFilePath)
	if err != nil {
		return Tracer{}, err
	}

	logger = logger.With("engine", "tracer")
	return Tracer{
		shouldTrace: cfg.ShouldTrace,
		Logger:      *logger,
	}, nil
}
