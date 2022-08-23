package kcl_plugin

import (
	"github.com/timandy/routine"
)

var ctxThreadLocal = routine.NewThreadLocal()

type PyPluginContext struct {
	PathList []string
	WorkDir  string
	Target   string
}

func NewPyPluginContext() *PyPluginContext {
	ctx := new(PyPluginContext)
	return ctx
}
