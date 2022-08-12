// Copyright 2021 The KCL Authors. All rights reserved.

pub const VERSION: &str = "0.4.2";
pub const CHECK_SUM: &str = "d6a5247bbfd8c157ac132105f4f57475";

pub fn get_full_version() -> String {
    format!("{}-{}", VERSION, CHECK_SUM)
}
