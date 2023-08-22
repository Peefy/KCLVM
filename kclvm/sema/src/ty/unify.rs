use std::collections::HashSet;

use super::{SchemaType, Type, TypeKind, TypeRef, UnsafeRef};

/// The type can be assigned to the expected type.
///
/// For security and performance considerations, dynamic dispatch of
/// types is not supported at this stage.
pub fn subsume(ty_lhs: TypeRef, ty_rhs: TypeRef, check_left_any: bool) -> bool {
    if (check_left_any && ty_lhs.is_any()) || (ty_rhs.is_any() || ty_lhs.is_none()) {
        true
    } else if ty_lhs.is_union() {
        let types = ty_lhs.union_types();
        types.iter().all(|ty| subsume(*ty, ty_rhs, false))
    } else if ty_rhs.is_union() {
        let types = ty_rhs.union_types();
        types.iter().any(|ty| subsume(ty_lhs, *ty, false))
    } else if ty_lhs.is_schema() {
        match &ty_rhs.kind {
            TypeKind::Schema(ty_rhs_schema) => {
                is_sub_schema_of(&ty_lhs.into_schema_type(), ty_rhs_schema)
            }
            _ => false,
        }
    } else if ty_lhs.is_int() && ty_rhs.is_float() {
        true
    } else if ty_lhs.is_number_multiplier() && ty_rhs.is_number_multiplier() {
        let ty_lhs = ty_lhs.into_number_multiplier();
        let ty_rhs = ty_rhs.into_number_multiplier();
        if ty_lhs.is_literal && ty_rhs.is_literal {
            ty_lhs.raw_value == ty_rhs.raw_value && ty_lhs.binary_suffix == ty_rhs.binary_suffix
        } else if ty_lhs.is_literal && !ty_rhs.is_literal {
            true
        } else {
            ty_lhs.is_literal || !ty_rhs.is_literal
        }
    } else if ty_lhs.is_primitive() && ty_rhs.is_primitive() {
        ty_lhs.kind == ty_rhs.kind
    } else if ty_lhs.is_literal() {
        if ty_rhs.is_literal() {
            ty_lhs.kind == ty_rhs.kind
        } else if ty_rhs.is_primitive() {
            // float_lit -> float
            // int_lit -> int
            // bool_lit -> bool
            // str_lit -> str
            // int_lit/bool_lit -> float
            if ty_rhs.is_float() && !ty_lhs.is_str() {
                true
            } else {
                ty_lhs.ty_str().contains(&ty_rhs.ty_str())
            }
        } else {
            false
        }
    } else if ty_lhs.is_list() && ty_rhs.is_list() {
        subsume(ty_lhs.list_item_ty(), ty_rhs.list_item_ty(), check_left_any)
    } else if ty_lhs.is_dict() && ty_rhs.is_dict() {
        let (ty_lhs_key, ty_lhs_val) = ty_lhs.dict_entry_ty();
        let (ty_rhs_key, ty_rhs_val) = ty_rhs.dict_entry_ty();
        subsume(ty_lhs_key, ty_rhs_key, check_left_any)
            && subsume(ty_lhs_val, ty_rhs_val, check_left_any)
    } else {
        equal(ty_lhs, ty_rhs)
    }
}

/// Are the two types exactly equal.
#[inline]
pub fn equal(ty_lhs: TypeRef, ty_rhs: TypeRef) -> bool {
    ty_lhs.kind == ty_rhs.kind
}

/// Whether the schema is sub schema of another schema.
pub fn is_sub_schema_of(schema_ty_lhs: &SchemaType, schema_ty_rhs: &SchemaType) -> bool {
    if schema_ty_lhs.ty_str_with_pkgpath() == schema_ty_rhs.ty_str_with_pkgpath() {
        true
    } else {
        match &schema_ty_lhs.base {
            Some(base) => is_sub_schema_of(base, schema_ty_rhs),
            None => false,
        }
    }
}

/// The type can be assigned to the expected type.
#[inline]
pub fn assignable_to(ty: TypeRef, expected_ty: TypeRef) -> bool {
    if !ty.is_assignable_type() {
        return false;
    }
    subsume(ty, expected_ty, true)
}

/// Whether `lhs_ty` is the upper bound of the `rhs_ty`
#[inline]
pub fn is_upper_bound(lhs_ty: TypeRef, rhs_ty: TypeRef) -> bool {
    subsume(rhs_ty, lhs_ty, false)
}

/// Whether the type list contains the `any` type.
#[inline]
pub fn has_any_type(types: &[TypeRef]) -> bool {
    types.iter().any(|ty| ty.is_any())
}

/// The sup function returns the minimum supremum of all types in an array of types.
#[inline]
pub fn sup(types: &[TypeRef]) -> TypeRef {
    r#typeof(types, true)
}

/// Typeof types
pub fn r#typeof(types: &[TypeRef], should_remove_sub_types: bool) -> TypeRef {
    // 1. Initialize an ordered set to store the type array
    let mut type_set: Vec<TypeRef> = vec![];
    // 2. Add the type array to the ordered set for sorting by the type id and de-duplication.
    add_types_to_type_set(&mut type_set, types);
    // 3. Remove sub types according to partial order relation rules e.g. sub schema types.
    if should_remove_sub_types {
        let mut remove_index_set = HashSet::new();
        for (i, source) in type_set.iter().enumerate() {
            for (j, target) in type_set.iter().enumerate() {
                if i != j && subsume(*source, *target, false) {
                    remove_index_set.insert(i);
                }
            }
        }
        let types: Vec<(usize, &TypeRef)> = type_set
            .iter()
            .enumerate()
            .filter(|(i, _)| !remove_index_set.contains(i))
            .collect();
        type_set = types.iter().map(|(_, ty)| *<&TypeRef>::clone(ty)).collect();
    }
    if type_set.is_empty() {
        UnsafeRef::new(Type::ANY)
    } else if type_set.len() == 1 {
        type_set[0]
    } else {
        Type::union_ref(&type_set)
    }
}

fn add_types_to_type_set(type_set: &mut Vec<TypeRef>, types: &[TypeRef]) {
    for ty in types {
        add_type_to_type_set(type_set, *ty);
    }
}

fn add_type_to_type_set(type_set: &mut Vec<TypeRef>, ty: TypeRef) {
    match &ty.kind {
        TypeKind::Union(types) => {
            add_types_to_type_set(type_set, types);
        }
        _ => {
            // Remove the bottom type.
            if !ty.is_void() && !type_set.contains(&ty) {
                type_set.push(ty)
            }
        }
    }
}
