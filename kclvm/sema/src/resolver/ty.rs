use crate::resolver::Resolver;
use crate::ty::parser::parse_type_str;
use crate::ty::{assignable_to, SchemaType, Type, TypeKind, TypeRef, UnsafeRef};
use indexmap::IndexMap;
use kclvm_ast::ast;
use kclvm_ast::pos::GetPos;
use kclvm_error::diagnostic::Range;
use kclvm_error::*;

use super::node::ResolvedResult;

fn ty_str_to_pkgpath(ty_str: &str) -> &str {
    let splits: Vec<&str> = ty_str.rsplitn(2, '.').collect();
    let len = splits.len();
    splits[len - 1]
}

pub fn ty_str_replace_pkgpath(ty_str: &str, pkgpath: &str) -> String {
    let pkgpath = format!("@{}", pkgpath);
    if ty_str.contains('.') && ty_str_to_pkgpath(ty_str) == pkgpath {
        ty_str.replacen(&format!("{}.", pkgpath), "", 1)
    } else {
        ty_str.to_string()
    }
}

impl<'ctx> Resolver<'ctx> {
    #[inline]
    pub fn any_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.any.clone()
    }
    #[inline]
    pub fn int_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.int.clone()
    }
    #[inline]
    pub fn float_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.float.clone()
    }
    #[inline]
    pub fn bool_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.bool.clone()
    }
    #[inline]
    pub fn str_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.str.clone()
    }
    #[inline]
    pub fn none_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.none.clone()
    }
    #[inline]
    pub fn void_ty(&self) -> TypeRef {
        self.ctx.ty_ctx.builtin_types.void.clone()
    }
    /// Parse the type string with the scope, if parse_ty returns a Named type(schema type or type alias),
    /// found it from the scope.
    pub fn parse_ty_with_scope(&mut self, ty: &ast::Type, range: Range) -> ResolvedResult {
        let ty: TypeRef = UnsafeRef::new(ty.clone().into());
        // If a named type, find it from scope to get the specific type
        let ret_ty = self.upgrade_named_ty_with_scope(ty.clone(), &range);
        self.add_type_alias(
            &ty.into_type_annotation_str(),
            &ret_ty.into_type_annotation_str(),
        );
        ret_ty
    }

    pub fn parse_ty_str_with_scope(&mut self, ty_str: &str, range: Range) -> ResolvedResult {
        let ty: TypeRef = parse_type_str(ty_str);
        // If a named type, find it from scope to get the specific type
        let ret_ty = self.upgrade_named_ty_with_scope(ty, &range);
        self.add_type_alias(ty_str, &ret_ty.into_type_annotation_str());
        ret_ty
    }

    /// The given expression must be the expected type.
    #[inline]
    pub fn must_be_type(&mut self, expr: &'ctx ast::NodeRef<ast::Expr>, expected_ty: TypeRef) {
        let ty = self.expr(expr);
        self.must_assignable_to(ty, expected_ty, expr.get_span_pos(), None);
    }

    /// Must assignable to the expected type.
    #[inline]
    pub fn must_assignable_to(
        &mut self,
        ty: TypeRef,
        expected_ty: TypeRef,
        range: Range,
        expected_pos: Option<Range>,
    ) {
        if !self.check_type(ty.clone(), expected_ty.clone(), &range) {
            let mut msgs = vec![Message {
                range,
                style: Style::LineAndColumn,
                message: format!("expected {}, got {}", expected_ty.ty_str(), ty.ty_str(),),
                note: None,
            }];

            if let Some(expected_pos) = expected_pos {
                msgs.push(Message {
                    range: expected_pos,
                    style: Style::LineAndColumn,
                    message: format!(
                        "variable is defined here, its type is {}, but got {}",
                        expected_ty.ty_str(),
                        ty.ty_str(),
                    ),
                    note: None,
                });
            }
            self.handler.add_error(ErrorKind::TypeError, &msgs);
        }
    }

    /// The check type main function, returns a boolean result.
    #[inline]
    pub fn check_type(&mut self, ty: TypeRef, expected_ty: TypeRef, range: &Range) -> bool {
        match (&ty.kind, &expected_ty.kind) {
            (TypeKind::List(item_ty), TypeKind::List(expected_item_ty)) => {
                self.check_type(item_ty.clone(), expected_item_ty.clone(), range)
            }
            (TypeKind::Dict(key_ty, val_ty), TypeKind::Dict(expected_key_ty, expected_val_ty)) => {
                self.check_type(key_ty.clone(), expected_key_ty.clone(), range)
                    && self.check_type(val_ty.clone(), expected_val_ty.clone(), range)
            }
            (TypeKind::Dict(key_ty, val_ty), TypeKind::Schema(schema_ty)) => {
                self.dict_assignable_to_schema(key_ty.clone(), val_ty.clone(), schema_ty, range)
            }
            (TypeKind::Union(types), _) => types
                .iter()
                .all(|ty| self.check_type(ty.clone(), expected_ty.clone(), range)),
            (_, TypeKind::Union(types)) => types
                .iter()
                .any(|expected_ty| self.check_type(ty.clone(), expected_ty.clone(), range)),
            _ => assignable_to(ty, expected_ty),
        }
    }

    /// Judge a dict can be converted to schema in compile time
    /// Do relaxed schema check key and value type check.
    pub fn dict_assignable_to_schema(
        &mut self,
        key_ty: TypeRef,
        val_ty: TypeRef,
        schema_ty: &SchemaType,
        range: &Range,
    ) -> bool {
        if let Some(index_signature) = &schema_ty.index_signature {
            if !assignable_to(val_ty, index_signature.val_ty) {
                self.handler.add_type_error(
                    &format!(
                        "expected schema index signature value type {}, got {}",
                        index_signature.val_ty.ty_str(),
                        val_ty.ty_str()
                    ),
                    range.clone(),
                );
            }
            if index_signature.any_other {
                return assignable_to(key_ty, index_signature.key_ty)
                    && assignable_to(val_ty, index_signature.val_ty);
            }
            true
        } else {
            true
        }
    }

    fn upgrade_named_ty_with_scope(&mut self, ty: TypeRef, range: &Range) -> ResolvedResult {
        match &ty.kind {
            TypeKind::List(item_ty) => {
                Type::list_ref(self.upgrade_named_ty_with_scope(*item_ty, range))
            }
            TypeKind::Dict(key_ty, val_ty) => Type::dict_ref(
                self.upgrade_named_ty_with_scope(*key_ty, range),
                self.upgrade_named_ty_with_scope(*val_ty, range),
            ),
            TypeKind::Union(types) => Type::union_ref(
                &types
                    .iter()
                    .map(|ty| self.upgrade_named_ty_with_scope(*ty, range))
                    .collect::<Vec<TypeRef>>(),
            ),
            TypeKind::Named(ty_str) => {
                let ty_str = ty_str_replace_pkgpath(ty_str, &self.ctx.pkgpath);
                let names: Vec<&str> = if ty_str.starts_with('@') {
                    let names: Vec<&str> = ty_str.rsplitn(2, '.').collect();
                    names.iter().rev().cloned().collect()
                } else {
                    ty_str.split('.').collect()
                };
                if names.is_empty() {
                    self.handler
                        .add_compile_error("missing type annotation", range.clone());
                    return self.any_ty();
                }
                let mut pkgpath = "".to_string();
                let name = names[0];
                if names.len() > 1 && !self.ctx.local_vars.contains(&name.to_string()) {
                    if let Some(mapping) = self.ctx.import_names.get(&self.ctx.filename) {
                        pkgpath = mapping
                            .get(name)
                            .map_or("".to_string(), |pkgpath| pkgpath.to_string());
                    }
                }
                self.ctx.l_value = false;
                self.resolve_var(
                    &names.iter().map(|n| n.to_string()).collect::<Vec<String>>(),
                    &pkgpath,
                    range.clone(),
                )
            }
            _ => ty,
        }
    }

    pub fn add_type_alias(&mut self, name: &str, alias: &str) {
        if alias.starts_with('@') {
            if name == &alias[1..] {
                return;
            }
        } else if name == alias {
            return;
        }
        match self.ctx.type_alias_mapping.get_mut(&self.ctx.pkgpath) {
            Some(mapping) => {
                mapping.insert(name.to_string(), alias.to_string());
            }
            None => {
                let mut mapping = IndexMap::default();
                mapping.insert(name.to_string(), alias.to_string());
                self.ctx
                    .type_alias_mapping
                    .insert(self.ctx.pkgpath.clone(), mapping);
            }
        }
    }
}
