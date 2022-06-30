package log

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-colorable"
	tmtLog "github.com/ethereum/go-ethereum/log"
)

func SetupLogger(level log.Lvl) {
	glogger := log.NewGlogHandler(log.StreamHandler(colorable.NewColorableStderr(), log.TerminalFormat(true)))
	glogger.Verbosity(level)

	log.Root().SetHandler(glogger)
}

//var _ tmtLog.Logger = (*Logger)(nil)

type Logger struct {
	keyvals []interface{}
}

func (l Logger) New(ctx ...interface{}) tmtLog.Logger {
	panic("implement me")
}

func (l Logger) GetHandler() tmtLog.Handler {
	panic("implement me")
}

func (l Logger) SetHandler(h tmtLog.Handler) {
	panic("implement me")
}

func (l Logger) Trace(msg string, ctx ...interface{}) {
	panic("implement me")
}

func (l Logger) Crit(msg string, ctx ...interface{}) {
	panic("implement me")
}

// Logger returns a new instance of an lightchain Logger.
// With() should  be called upon the returned instance to set context
//
// Logger is compatible with tmtLog.Logger.
func NewLogger() Logger {
	return Logger{keyvals: make([]interface{}, 0)}
}

// Debug proxies everything to the go-ethereum logging facilities
func (l Logger) Debug(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Debug(msg, ctx...)
}

// Info proxies everything to the go-ethereum logging facilities
func (l Logger) Info(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Info(msg, ctx...)
}

// Error proxies everything to the go-ethereum logging facilities
func (l Logger) Error(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Error(msg, ctx...)
}

// Error proxies everything to the go-ethereum logging facilities
func (l Logger) Warn(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Warn(msg, ctx...)
}


// With proxies everything to the go-ethereum logging facilities
func (l Logger) With(ctx ...interface{}) tmtLog.Logger {
	l.keyvals = append(l.keyvals, ctx...)
	return l
}
