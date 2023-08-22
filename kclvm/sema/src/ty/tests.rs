use super::*;

#[test]
fn test_sup() {
    let cases = vec![
        (vec![], UnsafeRef::new(Type::ANY)),
        (vec![UnsafeRef::new(Type::ANY)], UnsafeRef::new(Type::ANY)),
        (vec![UnsafeRef::new(Type::STR)], UnsafeRef::new(Type::STR)),
        (
            vec![UnsafeRef::new(Type::STR), UnsafeRef::new(Type::INT)],
            Type::union_ref(&[UnsafeRef::new(Type::STR), UnsafeRef::new(Type::INT)]),
        ),
        (
            vec![
                UnsafeRef::new(Type::BOOL),
                UnsafeRef::new(Type::bool_lit(true)),
            ],
            UnsafeRef::new(Type::BOOL),
        ),
        (
            vec![
                UnsafeRef::new(Type::str_lit("Blue")),
                UnsafeRef::new(Type::str_lit("Yellow")),
                UnsafeRef::new(Type::str_lit("Red")),
            ],
            Type::union_ref(&[
                UnsafeRef::new(Type::str_lit("Blue")),
                UnsafeRef::new(Type::str_lit("Yellow")),
                UnsafeRef::new(Type::str_lit("Red")),
            ]),
        ),
        (
            vec![
                Type::list_ref(Type::union_ref(&[
                    UnsafeRef::new(Type::int_lit(1)),
                    UnsafeRef::new(Type::int_lit(2)),
                ])),
                Type::list_ref(Type::union_ref(&[
                    UnsafeRef::new(Type::int_lit(3)),
                    UnsafeRef::new(Type::int_lit(4)),
                ])),
            ],
            Type::union_ref(&[
                Type::list_ref(Type::union_ref(&[
                    UnsafeRef::new(Type::int_lit(1)),
                    UnsafeRef::new(Type::int_lit(2)),
                ])),
                Type::list_ref(Type::union_ref(&[
                    UnsafeRef::new(Type::int_lit(3)),
                    UnsafeRef::new(Type::int_lit(4)),
                ])),
            ]),
        ),
        (
            vec![
                Type::union_ref(&[
                    UnsafeRef::new(Type::STR),
                    Type::dict_ref(UnsafeRef::new(Type::STR), UnsafeRef::new(Type::STR)),
                ]),
                Type::dict_ref(UnsafeRef::new(Type::ANY), UnsafeRef::new(Type::ANY)),
            ],
            Type::union_ref(&[
                UnsafeRef::new(Type::STR),
                Type::dict_ref(UnsafeRef::new(Type::ANY), UnsafeRef::new(Type::ANY)),
            ]),
        ),
    ];
    for (types, expected) in &cases {
        let got = sup(types);
        assert_eq!(got, *expected);
    }
}

#[test]
fn test_type_walker() {
    fn walk_fn(ty: &Type) -> TypeRef {
        if ty.is_int() {
            UnsafeRef::new(Type::STR)
        } else {
            UnsafeRef::new(ty.clone())
        }
    }
    let cases = [
        (UnsafeRef::new(Type::ANY), UnsafeRef::new(Type::ANY)),
        (UnsafeRef::new(Type::INT), UnsafeRef::new(Type::STR)),
        (UnsafeRef::new(Type::STR), UnsafeRef::new(Type::STR)),
        (
            Type::list_ref(UnsafeRef::new(Type::INT)),
            Type::list_ref(UnsafeRef::new(Type::STR)),
        ),
        (
            Type::union_ref(&[UnsafeRef::new(Type::INT), UnsafeRef::new(Type::STR)]),
            Type::union_ref(&[UnsafeRef::new(Type::STR), UnsafeRef::new(Type::STR)]),
        ),
        (
            Type::union_ref(&[
                UnsafeRef::new(Type::INT),
                UnsafeRef::new(Type::STR),
                Type::union_ref(&[UnsafeRef::new(Type::INT), UnsafeRef::new(Type::STR)]),
            ]),
            Type::union_ref(&[
                UnsafeRef::new(Type::STR),
                UnsafeRef::new(Type::STR),
                Type::union_ref(&[UnsafeRef::new(Type::STR), UnsafeRef::new(Type::STR)]),
            ]),
        ),
        (
            Type::dict_ref(UnsafeRef::new(Type::INT), UnsafeRef::new(Type::INT)),
            Type::dict_ref(UnsafeRef::new(Type::STR), UnsafeRef::new(Type::STR)),
        ),
    ];
    for (ty, expected) in cases {
        assert_eq!(
            walker::walk_type(&ty, walk_fn),
            expected,
            "Type test failed: {}",
            ty.ty_str()
        );
    }
}
