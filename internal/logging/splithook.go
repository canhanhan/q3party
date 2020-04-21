package logging

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type SplitHook struct{}

func (h *SplitHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel}
}

func (h *SplitHook) Fire(entry *logrus.Entry) error {
	msg := fmt.Sprintf("%s: %s (%s)\n", entry.Level, entry.Message, entry.Time)
	if entry.Level == logrus.InfoLevel {
		os.Stdout.WriteString(msg)
	} else {
		os.Stderr.WriteString(msg)
	}

	return nil
}
