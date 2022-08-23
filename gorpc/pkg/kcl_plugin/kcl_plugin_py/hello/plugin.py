# Copyright 2020 The KCL Authors. All rights reserved.

INFO = {
    'name': 'hello',
    'describe': 'hello doc',
    'long_describe': 'long describe',
    'version': '0.0.1',
}


global_int: int = 0


def set_global_int(v: int):
    global global_int
    global_int = v


def get_global_int() -> int:
    return global_int


def say_hello(msg: str):
    print('hello.say_hello:', msg)
    return None


def add(a: int, b: int) -> int:
    """add two numbers, and return result"""
    return a + b


def tolower(s: str) -> str:
    return s.lower()


def update_dict(d: dict, key: str, value: str) -> dict:
    d[key] = value
    return d


def list_append(l: list, *values) -> list:
    for v in values:
        l.append(v)
    return l


def foo(a, b, *, x, **values):
    print(a, b, x, values)
    return {'a': a, 'b': b, 'x': x, **values}
