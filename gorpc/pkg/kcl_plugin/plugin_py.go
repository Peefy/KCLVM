// Copyright 2022 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package kcl_plugin

/*
typedef struct PyObject PyObject;
*/
import "C"
import (
	"fmt"
)

func py_callPluginMethod(method, args_json, kwargs_json string) string {
	result := ""
	PyGILEntry(func() {
		ctx := ctxThreadLocal.Get()
		switch ctx := ctx.(type) {
		case *PyPluginContext:
			SaveKclvmPyConfig()
			SetKclvmPyConfig(ctx.PathList, ctx.WorkDir, ctx.Target)
			py_method := PyUnicodeFromString(method)
			py_args_json := PyUnicodeFromString(args_json)
			py_kwargs_json := PyUnicodeFromString(kwargs_json)
			py_result := CallPyFunc(kclvmPyPluginModule, "_call_py_method", map[string]*C.PyObject{}, py_method, py_args_json, py_kwargs_json)
			result = PyUnicodeAsUTF8(py_result)
			RecoverKclvmPyConfig()
		default:
			result = JSONError(fmt.Errorf("ctxThreadLocal can't get PyPluginContext: %v", ctx))
		}
	})
	return result
}
