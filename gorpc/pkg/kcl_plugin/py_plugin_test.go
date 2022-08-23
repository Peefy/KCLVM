// Copyright 2022 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package kcl_plugin

import (
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPyPlugin(t *testing.T) {
	_, srcPath, _, _ := runtime.Caller(0)
	CppSetEnv("KCL_PLUGINS_ROOT", srcPath[:strings.LastIndex(srcPath, "/")]+"/kcl_plugin_py", true)
	ctxThreadLocal.Set(NewPyPluginContext())
	py_callPluginMethod("hello.say_hello", "[\"kclvm-go\"]", "")
	result := py_callPluginMethod("hello.add", "[1,2]", "")
	assert.Equal(t, "3", result)
	result = py_callPluginMethod("hello.tolower", "[\"KCL\"]", "")
	assert.Equal(t, "\"kcl\"", result)
	result = py_callPluginMethod("hello.get_global_int", "[]", "")
	assert.Equal(t, "0", result)
	py_callPluginMethod("hello.set_global_int", "[9973]", "")
	result = py_callPluginMethod("hello.get_global_int", "[]", "")
	assert.Equal(t, "9973", result)
	result = py_callPluginMethod("hello.foo", "[\"aaa\",\"bbb\"]", "{\"x\":123,\"y\":234,\"acbd\":1234}")
	assert.Equal(t, "{\"a\": \"aaa\", \"b\": \"bbb\", \"x\": 123, \"y\": 234, \"acbd\": 1234}", result)
	result = py_callPluginMethod("hello.foo", "[\"aaa\",\"bbb\"]", "{\"x\":123}")
	assert.Equal(t, "{\"a\": \"aaa\", \"b\": \"bbb\", \"x\": 123}", result)
}

func TestPyPluginInMultiGoRoutine(t *testing.T) {
	_, srcPath, _, _ := runtime.Caller(0)
	CppSetEnv("KCL_PLUGINS_ROOT", srcPath[:strings.LastIndex(srcPath, "/")]+"/kcl_plugin_py", true)
	n := 1000
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			m := i
			ctxThreadLocal.Set(NewPyPluginContext())
			result := py_callPluginMethod("hello.add", "[1,2]", "")
			println(m)
			assert.Equal(t, "3", result)
			wg.Done()
		}()
	}
	wg.Wait()
}
