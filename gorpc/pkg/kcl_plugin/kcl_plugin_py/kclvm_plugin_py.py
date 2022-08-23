import typing
import json
import inspect
import sys
import os

import kclvm.config
import kclvm.compiler.vfs as vfs
import kclvm.kcl.info as kcl_info
import kclvm.compiler.extension.plugin.plugin as kcl_plugin

_plugin_dict = {}

KCLVM_TARGET_ENV_KEY = "KCLVM_TARGET"

saved_input_file = None
saved_current_path = ""
saved_is_target_native  = True
saved_is_target_wasm  = False

def _call_py_method(name: str, args_json: str, kwargs_json: str) -> str:
        try:
            return _call_py_method_unsafe(name, args_json, kwargs_json)
        except Exception as e:
            return json.dumps({"__kcl_PanicInfo__": f"{e}"})

def _call_py_method_unsafe(
        name: str, args_json: str, kwargs_json: str
    ) -> str:
        dotIdx = name.rfind(".")
        if dotIdx < 0:
            return ""
        
        modulePath = name[:dotIdx]
        mathodName = name[dotIdx + 1 :]

        plugin_name = modulePath[modulePath.rfind(".") + 1 :]

        module = _get_plugin(plugin_name)
        mathodFunc = None
        for func_name, func in inspect.getmembers(module):
            if func_name == kcl_info.demangle(mathodName):
                mathodFunc = func
                break
        args = []
        kwargs = {}

        if args_json:
            args = json.loads(args_json)
            if not isinstance(args, list):
                return ""
        if kwargs_json:
            kwargs = json.loads(kwargs_json)
            if not isinstance(kwargs, dict):
                return ""

        result = mathodFunc(*args, **kwargs)
        sys.stdout.flush()
        return json.dumps(result)

def _get_plugin(plugin_name: str) -> typing.Optional[any]:
    if plugin_name in _plugin_dict:
        return _plugin_dict[plugin_name]
    
    module = kcl_plugin.get_plugin(plugin_name)
    _plugin_dict[plugin_name] = module
    return module

def _get_target(path_list: typing.List[str]) -> str:
    root = vfs.MustGetPkgRoot(path_list)
    modfile = vfs.LoadModFile(root)
    return (modfile.build.target or os.getenv(KCLVM_TARGET_ENV_KEY) or "").lower()

def _set_kclvm_config(path_list,work_dir: str,target: str):
    kclvm.config.input_file = path_list
    kclvm.config.current_path = work_dir
    kclvm.config.is_target_native = target == "native"
    kclvm.config.is_target_wasm = target == "wasm"

def _save_kclvm_config():
    saved_input_file = kclvm.config.input_file 
    saved_current_path = kclvm.config.current_path
    saved_is_target_native  = kclvm.config.is_target_native
    saved_is_target_wasm  = kclvm.config.is_target_wasm

def _recover_kclvm_config():
    kclvm.config.input_file=saved_input_file
    kclvm.config.current_path=saved_current_path
    kclvm.config.is_target_native=saved_is_target_native
    kclvm.config.is_target_wasm=saved_is_target_wasm

def hello(name :str) -> str:
    return "hello plugin : " +name