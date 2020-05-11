package logging

import "go.uber.org/zap"

const callerSkip = 1

// NewLogger creates and returns a Zap Logger
func NewLogger() *zap.Logger {
	l, err := zap.NewDevelopment(zap.AddCallerSkip(callerSkip))
	if err != nil {
		l = zap.L()
	}
	return l
}
