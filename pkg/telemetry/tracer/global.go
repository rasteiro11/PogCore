package tracer

import (
	"sync"

	"github.com/rasteiro11/PogCore/pkg/logger"
)

var (
	global Tracer
	once   sync.Once
)

func Instance() Tracer {
	once.Do(func() {
		if global == nil {
			global = &noopTracer{}
		}
	})

	return global
}

func SetGlobal(t Tracer) {
	if t == nil {
		logger.Global().Warn("[tracer.SetGlobal] trace is null, setting a noop provider to prevent errors")
		global = &noopTracer{}
		return
	}

	global = t
}
