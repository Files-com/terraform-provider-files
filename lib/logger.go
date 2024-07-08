package lib

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Logger struct {
	Ctx context.Context
}

func (l Logger) Printf(msg string, args ...interface{}) {
	l.Debug(msg, args...)
}
func (l Logger) Error(msg string, args ...interface{}) {
	tflog.Error(l.Ctx, fmt.Sprintf(msg, args...))
}
func (l Logger) Info(msg string, args ...interface{}) {
	tflog.Info(l.Ctx, fmt.Sprintf(msg, args...))
}
func (l Logger) Debug(msg string, args ...interface{}) {
	tflog.Debug(l.Ctx, fmt.Sprintf(msg, args...))
}
func (l Logger) Warn(msg string, args ...interface{}) {
	tflog.Warn(l.Ctx, fmt.Sprintf(msg, args...))
}
