// Copyright 2022 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package kcl_plugin

/*

#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>

typedef
    enum {PyGILState_LOCKED, PyGILState_UNLOCKED}
        Kclvm_PyGILState_STATE;
typedef ssize_t Py_ssize_t;
typedef struct PyCompilerFlags PyCompilerFlags;
typedef struct PyObject PyObject;
typedef struct PyThreadState PyThreadState;
int Kclvm_Py_Main(void * f,int argc,wchar_t **argv){
   int (*py_main)(int,wchar_t **);
   py_main =(int (*)(int,wchar_t **))f;
   return py_main(argc,argv);
}

int Kclvm_PyRun_AnyFileFlags(void * f,FILE * file, const char * arg, PyCompilerFlags * flag){
   int (*py_run_file)(FILE *, const char *, PyCompilerFlags *);
   py_run_file = (int (*)(FILE *, const char *, PyCompilerFlags *))f;
   return py_run_file(file,arg,flag);
}

void  Kclvm_Py_SetPythonHome(void * f,const wchar_t * home){
	void (*set_home)(const wchar_t *);
	set_home = (void (*)(const wchar_t *))f;
	set_home(home);
}

wchar_t * Kclvm_Py_DecodeLocale(void * f,const char *arg,size_t *size){
   wchar_t *(*py_decode)(const char*,size_t*);
   py_decode =(wchar_t *(*)(const char*,size_t*))f;
   return py_decode(arg,size);
}

void Kclvm_PyMem_RawFree(void * f ,void *ptr){
   void (*py_free)(void *);
   py_free = (void (*)(void *))f;
   return py_free(ptr);
}

void Kclvm_Py_Initialize(void * f){
   void (*py_init)(void);
   py_init =  (void (*)(void))f;
   py_init();
}

void Kclvm_Py_InitializeEx(void * f,int i) {
   void (*py_init)(int);
   py_init = (void (*)(int))f;
   py_init(i);
}

void Kclvm_PyEval_InitThreads(void * f){
   void (*thread_init)(void);
   thread_init =  (void (*)(void))f;
   thread_init();
}

PyThreadState * Kclvm_PyEval_SaveThread(void * f){
   PyThreadState* (*save_thread)(void);
   save_thread =  (PyThreadState *(*)(void))f;
   return save_thread();
}

Kclvm_PyGILState_STATE Kclvm_PyGILState_Ensure(void * f){
   Kclvm_PyGILState_STATE (*state_ensure)(void);
   state_ensure =  (Kclvm_PyGILState_STATE (*)(void))f;
   return state_ensure();
}

void Kclvm_PyGILState_Release(void * f,Kclvm_PyGILState_STATE state){
   void (*release_state)(Kclvm_PyGILState_STATE);
   release_state =  (void (*)(Kclvm_PyGILState_STATE))f;
   release_state(state);
}

int Kclvm_PyGILState_Check(void * f){
	int (*check_state)(void);
	check_state = (int (*)(void))f;
	return check_state();
}

PyObject * Kclvm_PyImport_ImportModule(void * f,const char *name){
	PyObject * (*import_mod)(const char *);
	import_mod = (PyObject * (*)(const char *name))f;
	return import_mod(name);
}

PyObject * Kclvm_PyImport_ReloadModule(void * f ,PyObject *m){
	PyObject * (*reload)(PyObject *);
	reload = (PyObject * (*)(PyObject *))f;
	return reload(m);
}

PyObject * Kclvm_PyObject_GetAttrString(void * f,PyObject * obj, const char * name){
	PyObject * (*get_attr)(PyObject *, const char *);
	get_attr = (PyObject * (*)(PyObject *, const char *))f;
	return get_attr(obj,name);
}

int Kclvm_PyObject_SetAttrString(void * f,PyObject * obj, const char * attr_name , PyObject * v){
	int (*set_attr)(PyObject *, const char *, PyObject *);
	set_attr = (int (*)(PyObject *, const char *, PyObject *))f;
	return set_attr(obj,attr_name,v);
}

PyObject * Kclvm_PyObject_Call(void *f,PyObject *callable, PyObject *args, PyObject *kwargs){
	PyObject * (*call)(PyObject *, PyObject *, PyObject *);
	call = (PyObject * (*)(PyObject *, PyObject *, PyObject *))f;
	return call(callable,args,kwargs);
}

PyObject * Kclvm_PyObject_CallMethodObjArgs(void * f,PyObject *obj, PyObject *name,...){
	PyObject * (*call_method)(PyObject *, PyObject *,...);
	call_method = (PyObject * (*)(PyObject *, PyObject *,...))f;
	va_list arg_list;
	va_start(arg_list, name);
	PyObject* result = call_method(obj,name);
	va_end(arg_list);
	return result;
}

PyObject * Kclvm_PyObject_CallMethodNoArgs(void *f, PyObject * obj , PyObject *name){
	PyObject * (*call_method)(PyObject *, PyObject *,...);
	call_method = (PyObject * (*)(PyObject *, PyObject *,...))f;
	return call_method(obj,name);
}

PyObject * Kclvm_PyObject_CallMethodOneArgs(void *f, PyObject * obj , PyObject *name,PyObject * arg){
	PyObject * (*call_method)(PyObject *, PyObject *,...);
	call_method = (PyObject * (*)(PyObject *, PyObject *,...))f;
	return call_method(obj,name,arg);
}

PyObject *  Kclvm_PyObject_Repr(void * f,PyObject *obj){
	PyObject *  (*repr)(PyObject *);
	repr = (PyObject *  (*)(PyObject *))f;
	return repr(obj);
}

PyObject *Kclvm_PySys_GetObject(void * f,const char *name){
	PyObject *(*get_obj)(const char *name);
	get_obj = (PyObject *(*)(const char *name))f;
	return get_obj(name);
}

PyObject* Kclvm_PyUnicode_FromString(void * f,const char *u){
	PyObject* (*py_unicode)(const char *);
	py_unicode = (PyObject* (*)(const char *))f;
	return py_unicode(u);
}

const char * Kclvm_PyUnicode_AsUTF8(void * f,PyObject *unicode){
	const char * (*as_utf8)(PyObject *);
	as_utf8 = (const char * (*)(PyObject *))f;
	return as_utf8(unicode);
}

void Kclvm_Py_DecRef(void * f,PyObject * obj){
	void (*dec_ref)(PyObject *);
	dec_ref = (void (*)(PyObject *))f;
	dec_ref(obj);
}

PyObject * Kclvm_PyList_New(void * f,Py_ssize_t size){
	PyObject * (*new_list)(Py_ssize_t);
	new_list = (PyObject * (*)(Py_ssize_t))f;
	return new_list(size);
}

int Kclvm_PyList_SetItem(void * f,PyObject * obj, Py_ssize_t index, PyObject * item){
	int (*set_item)(PyObject *, Py_ssize_t, PyObject *);
	set_item = (int (*)(PyObject *, Py_ssize_t, PyObject *))f;
	return set_item(obj,index,item);
}

int Kclvm_PyList_Append(void * f,PyObject * obj, PyObject * item) {
	int (*append)(PyObject *, PyObject *);
	append = (int (*)(PyObject *, PyObject *))f;
	return append(obj,item);
}

int Kclvm_PyList_Insert(void * f,PyObject * obj, Py_ssize_t index, PyObject * item) {
	int (*insert)(PyObject *, Py_ssize_t, PyObject *);
	insert = (int (*)(PyObject *, Py_ssize_t, PyObject *))f;
	return insert(obj,index,item);
}

PyObject * Kclvm_PyTuple_New(void * f,Py_ssize_t size){
	PyObject * (*new_tuple)(Py_ssize_t);
	new_tuple = (PyObject * (*)(Py_ssize_t))f;
	return new_tuple(size);
}

int Kclvm_PyTuple_SetItem(void * f,PyObject * obj, Py_ssize_t index, PyObject * item){
	int (*set_item)(PyObject *, Py_ssize_t, PyObject *);
	set_item = (int (*)(PyObject *, Py_ssize_t, PyObject *))f;
	return set_item(obj,index,item);
}

PyObject * Kclvm_PyDict_New(void * f){
	PyObject * (*new_dict) (void);
	new_dict = (PyObject * (*) (void))f;
	return new_dict();
}


int Kclvm_PyDict_SetItem(void * f,PyObject *mp, PyObject *key, PyObject *item) {
	int (*set_item)(PyObject *, PyObject *, PyObject *);
	set_item = (int (*)(PyObject *, PyObject *, PyObject *))f;
	return set_item(mp,key,item);
}

void Kclvm_PyErr_Clear(void * f){
   void (*clear)(void);
   clear =  (void (*)(void))f;
   clear();
}

*/
import "C"

import (
	"fmt"
	"unsafe"
)

var pythonHomePtr *C.wchar_t

// see in https://docs.python.org/3/c-api/veryhigh.html#c.Py_Main
func PyMain(args []string) C.int {
	f := "Py_Main"

	funcPtr, _ := pyLib.GetSymbolPointer(f)
	argc := C.int(len(args))
	argv := make([]*C.wchar_t, argc, argc)
	for i, arg := range args {
		carg := C.CString(arg)
		defer C.free(unsafe.Pointer(carg))

		warg := PyDecodeLocale(carg, nil)
		if warg == nil {
			return -1
		}
		defer PyMemRawFree(unsafe.Pointer(warg))
		argv[i] = warg
	}

	return C.Kclvm_Py_Main(funcPtr, argc, (**C.wchar_t)(unsafe.Pointer(&argv[0])))
}

// see in https://docs.python.org/3/c-api/veryhigh.html#c.PyRun_AnyFile
func PyRunAnyFile(filename string) (int, error) {
	f := "PyRun_AnyFileFlags"
	funcPtr, _ := pyLib.GetSymbolPointer(f)

	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))

	cfile, err := C.fopen(cfilename, mode)

	if err != nil {
		return -1, fmt.Errorf("fail to open '%s': %s", filename, err)
	}

	defer C.fclose(cfile)

	return int(C.Kclvm_PyRun_AnyFileFlags(funcPtr, cfile, cfilename, nil)), nil
}

// see in https://docs.python.org/3/c-api/init.html#c.Py_SetPythonHome
func PySetPythonHome(home string) error {
	f := "Py_SetPythonHome"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	chome := C.CString(home)
	defer C.free(unsafe.Pointer(chome))

	newHome := PyDecodeLocale(chome, nil)
	if newHome == nil {
		return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", home)
	}
	C.Kclvm_Py_SetPythonHome(funcPtr, newHome)

	if pythonHomePtr != nil {
		PyMemRawFree(unsafe.Pointer(pythonHomePtr))
	}
	pythonHomePtr = newHome

	return nil
}

// see in https://docs.python.org/3/c-api/sys.html#c.Py_DecodeLocale
func PyDecodeLocale(arg *C.char, size *C.size_t) *C.wchar_t {
	f := "Py_DecodeLocale"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	return C.Kclvm_Py_DecodeLocale(funcPtr, arg, size)
}

// see in https://docs.python.org/3/c-api/memory.html#c.PyMem_RawFree
func PyMemRawFree(ptr unsafe.Pointer) {
	f := "PyMem_RawFree"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	C.Kclvm_PyMem_RawFree(funcPtr, ptr)

}

// see in https://docs.python.org/3/c-api/init.html#c.Py_Initialize
func PyInitialize() {
	f := "Py_Initialize"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	C.Kclvm_Py_Initialize(funcPtr)
}

// see in https://docs.python.org/3/c-api/init.html#c.Py_InitializeEx
func PyInitializeEx(i bool) {
	f := "Py_InitializeEx"

	funcPtr, _ := pyLib.GetSymbolPointer(f)
	if i {
		C.Kclvm_Py_InitializeEx(funcPtr, C.int(1))
	} else {
		C.Kclvm_Py_InitializeEx(funcPtr, C.int(0))
	}
}

// see in https://docs.python.org/3/c-api/import.html#c.PyImport_ImportModule
func PyImportImportModule(name string) *C.PyObject {
	f := "PyImport_ImportModule"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return C.Kclvm_PyImport_ImportModule(funcPtr, cname)

}

// see in https://docs.python.org/3/c-api/import.html#c.PyImport_ReloadModule
func PyImportReloadModule(m *C.PyObject) *C.PyObject {
	f := "PyImport_ReloadModule"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyImport_ReloadModule(funcPtr, m)
}

// see in https://docs.python.org/3/c-api/object.html#c.PyObject_GetAttrString
func PyObjectGetAttrString(obj *C.PyObject, name string) *C.PyObject {
	f := "PyObject_GetAttrString"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.Kclvm_PyObject_GetAttrString(funcPtr, obj, cname)
}

// see in https://docs.python.org/3/c-api/object.html#c.PyObject_SetAttrString
func PyObjectSetAttrString(obj *C.PyObject, attr_name string, v *C.PyObject) int {
	f := "PyObject_SetAttrString"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return int(C.Kclvm_PyObject_SetAttrString(funcPtr, obj, cattr_name, v))
}

// see in https://docs.python.org/3/c-api/object.html#c.PyObject_Call
func PyObjectCall(callable *C.PyObject, args *C.PyObject, kwargs *C.PyObject) *C.PyObject {
	f := "PyObject_Call"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyObject_Call(funcPtr, callable, args, kwargs)
}

// see in https://docs.python.org/3/c-api/call.html#c.PyObject_CallMethodNoArgs
func PyObjectCallMethodNoArgs(obj *C.PyObject, name *C.PyObject) *C.PyObject {
	f := "PyObject_CallMethodObjArgs"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyObject_CallMethodNoArgs(funcPtr, obj, name)
}

// see in https://docs.python.org/3/c-api/call.html#c.PyObject_CallMethodOneArgs
func PyObjectCallMethodOneArgs(obj *C.PyObject, name *C.PyObject, arg *C.PyObject) *C.PyObject {
	f := "PyObject_CallMethodObjArgs"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyObject_CallMethodOneArgs(funcPtr, obj, name, arg)
}

// see in https://docs.python.org/3/c-api/object.html#c.PyObject_Repr
func PyObjectRepr(obj *C.PyObject) *C.PyObject {
	f := "PyObject_Repr"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	return C.Kclvm_PyObject_Repr(funcPtr, obj)
}

// see in https://docs.python.org/3/c-api/sys.html#c.PySys_GetObject
func PySysGetObject(name string) *C.PyObject {
	f := "PySys_GetObject"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.Kclvm_PySys_GetObject(funcPtr, cName)
}

// see in https://docs.python.org/3/c-api/unicode.html#c.PyUnicode_FromString
func PyUnicodeFromString(u string) *C.PyObject {
	f := "PyUnicode_FromString"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	cu := C.CString(u)
	defer C.free(unsafe.Pointer(cu))
	return C.Kclvm_PyUnicode_FromString(funcPtr, cu)
}

// see in https://docs.python.org/3/c-api/unicode.html#c.PyUnicode_AsUTF8
func PyUnicodeAsUTF8(obj *C.PyObject) string {
	f := "PyUnicode_AsUTF8"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	return C.GoString(C.Kclvm_PyUnicode_AsUTF8(funcPtr, obj))
}

// see in https://docs.python.org/3/c-api/refcounting.html#c.Py_DecRef
func PyDecRef(obj *C.PyObject) {
	f := "Py_DecRef"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	C.Kclvm_Py_DecRef(funcPtr, obj)
}

// see in https://docs.python.org/3/c-api/list.html#c.PyList_New
func PyListNew(len int) *C.PyObject {
	f := "PyList_New"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyList_New(funcPtr, C.Py_ssize_t(len))
}

// see in https://docs.python.org/3/c-api/list.html#c.PyList_SetItem
func PyListSetItem(obj *C.PyObject, index int, item *C.PyObject) int {

	f := "PyList_SetItem"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return int(C.Kclvm_PyList_SetItem(funcPtr, obj, C.Py_ssize_t(index), item))
}

// see in https://docs.python.org/3/c-api/list.html#c.PyList_Append
func PyListAppend(obj, item *C.PyObject) int {
	f := "PyList_Append"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	return int(C.Kclvm_PyList_Append(funcPtr, obj, item))
}

// see in https://docs.python.org/3/c-api/list.html#c.PyList_Insert
func PyListInsert(p *C.PyObject, index int, item *C.PyObject) int {
	f := "PyList_Insert"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return int(C.Kclvm_PyList_Insert(funcPtr, p, C.Py_ssize_t(index), item))
}

// see in https://docs.python.org/3/c-api/tuple.html#c.PyTuple_New
func PyTupleNew(len int) *C.PyObject {
	f := "PyTuple_New"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyTuple_New(funcPtr, C.Py_ssize_t(len))
}

// see in https://docs.python.org/3/c-api/tuple.html#c.PyTuple_SetItem
func PyTupleSetItem(obj *C.PyObject, index int, item *C.PyObject) int {

	f := "PyTuple_SetItem"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return int(C.Kclvm_PyTuple_SetItem(funcPtr, obj, C.Py_ssize_t(index), item))
}

// see in https://docs.python.org/3/c-api/dict.html#c.PyDict_New
func PyDictNew() *C.PyObject {
	f := "PyDict_New"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	return C.Kclvm_PyDict_New(funcPtr)
}

// see in https://docs.python.org/3/c-api/dict.html#c.PyDict_SetItem
func PyDictSetItem(mp, key, item *C.PyObject) int {
	f := "PyDict_SetItem"
	funcPtr, _ := pyLib.GetSymbolPointer(f)
	return int(C.Kclvm_PyDict_SetItem(funcPtr, mp, key, item))
}

// see in https://docs.python.org/3/c-api/exceptions.html#c.PyErr_Clear
func PyErrClear() {
	f := "PyErr_Clear"

	funcPtr, _ := pyLib.GetSymbolPointer(f)

	C.Kclvm_PyErr_Clear(funcPtr)
}

func CppSetEnv(key, value string, overwrite bool) int {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	if overwrite {
		return int(C.setenv(cKey, cValue, C.int(1)))
	} else {
		return int(C.setenv(cKey, cValue, C.int(0)))
	}
}
