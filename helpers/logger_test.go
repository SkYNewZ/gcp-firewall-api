package helpers

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInitLogger(t *testing.T) {
	InitLogger()

	if logrus.GetLevel() != logrus.DebugLevel {
		t.Errorf("Log level should be logrus.DebugLevel. Got '%v' want '%v'", logrus.GetLevel(), logrus.DebugLevel)
	}
}
