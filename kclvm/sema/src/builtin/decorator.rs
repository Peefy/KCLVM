use indexmap::IndexMap;
use once_cell::sync::Lazy;

use crate::ty::{Parameter, Type, UnsafeRef};

macro_rules! register_decorator {
    ($($name:ident => $ty:expr)*) => (
        // Builtin decorator map.
        pub const BUILTIN_DECORATORS: Lazy<IndexMap<String, Type>> = Lazy::new(|| {
            let mut builtin_mapping = IndexMap::default();
            $( builtin_mapping.insert(stringify!($name).to_string(), $ty); )*
            builtin_mapping
        });
        pub static DECORATOR_NAMES: &[&str] = &[
            $( stringify!($name), )*
        ];
    )
}

register_decorator! {
    deprecated => Type::function(
        None,
        UnsafeRef::new(Type::ANY),
        &[
            Parameter {
                name: "version".to_string(),
                ty: UnsafeRef::new(Type::STR),
                has_default: true,
            },
            Parameter {
                name: "reason".to_string(),
                ty: UnsafeRef::new(Type::STR),
                has_default: true,
            },
            Parameter {
                name: "strict".to_string(),
                ty: UnsafeRef::new(Type::BOOL),
                has_default: true,
            },
        ],
        r#"This decorator is used to get the deprecation message according to the wrapped key-value pair.

        Examples
        --------
        @deprecated(version="v1.16", reason="The age attribute was deprecated", strict=True)
        schema Person:
            name: str
            age: int
        "#,
        false,
        None,
    )
    info => Type::function(
        None,
        UnsafeRef::new(Type::ANY),
        &[],
        r#"Info decorator is used to mark some compile-time information for external API queries

        Examples
        --------
        @info(message="User message")
        schema Person:
            name: str
            age: int
        "#,
        true,
        Some(0),
    )
}
