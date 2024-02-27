//! Copyright The KCL Authors. All rights reserved.

use std::error;
use std::fmt::{self, Debug};

pub(crate) const COMPILE_ERROR_MSG: &str = "Compile error";
pub(crate) const INTERNAL_ERROR_MSG: &str = "Internal error, please report a bug to us";

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum KCLErrorType {
    Compile,
    Runtime,
}

#[derive(Debug, Clone)]
pub struct KCLError {
    pub message: String,
    pub ty: KCLErrorType,
}

impl fmt::Display for KCLError {
    fn fmt(&self, f: &mut fmt::Formatter) -> Result<(), fmt::Error> {
        write!(
            f,
            "{}: {}",
            match self.ty {
                KCLErrorType::Compile => "compile error",
                KCLErrorType::Runtime => "runtime error",
            },
            self.message
        )
    }
}

impl error::Error for KCLError {}

impl Default for KCLError {
    fn default() -> Self {
        Self {
            message: Default::default(),
            ty: KCLErrorType::Compile,
        }
    }
}

impl KCLError {
    pub fn new(msg: &str) -> Self {
        Self {
            message: msg.to_string(),
            ty: KCLErrorType::Compile,
        }
    }
}
