package logging

import "go.uber.org/zap"

// NewLogger creates and returns a Zap Logger
func NewLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		l = zap.L()
	}
	return l
}
