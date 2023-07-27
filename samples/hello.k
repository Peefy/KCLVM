import kcl_plugin.hello

cc = hello.add(1, 2)
name = "kcl"
age = 1

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {}
x1 = Person {
    age = 101
}
