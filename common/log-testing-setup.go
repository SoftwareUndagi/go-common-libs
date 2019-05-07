package common

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

// LogCapturer reroutes testing.T log output
type LogCapturer interface {
	Release()
}

type logCapturer struct {
	*testing.T
	origOut io.Writer
}

func (tl logCapturer) Write(p []byte) (n int, err error) {
	tl.Logf((string)(p))
	return len(p), nil
}

func (tl logCapturer) Release() {
	logrus.SetOutput(tl.origOut)
}

// CaptureLog redirects logrus output to testing.Log
func CaptureLog(t *testing.T) LogCapturer {
	lc := logCapturer{T: t, origOut: logrus.StandardLogger().Out}
	logrus.SetOutput(lc)
	if !testing.Verbose() {
		logrus.SetOutput(lc)
	}
	return &lc
}
