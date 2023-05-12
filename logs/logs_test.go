package logs

import (
	"testing"
)

func TestDebug(t *testing.T) {
	Debug("debug test")
}

func TestDebugf(t *testing.T) {
	Debugf("%s%s", "1", "2")
}

func TestInfo(t *testing.T) {
	Info("info test")
}

func TestInfof(t *testing.T) {
	Infof("info: %s%s", "1", "2")
}

func TestError(t *testing.T) {
	Error("error test")
}

func TestErrorf(t *testing.T) {
	Errorf("error: %s%s", "1", "2")
}
