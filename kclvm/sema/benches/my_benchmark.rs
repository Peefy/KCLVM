use criterion::{criterion_group, criterion_main, Criterion};
use kclvm_sema::ty::*;

pub fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("sup", |b| {
        b.iter(|| {
            let types = vec![
                UnsafeRef::new(Type::int_lit(1)),
                UnsafeRef::new(Type::INT),
                UnsafeRef::new(Type::union(&[
                    UnsafeRef::new(Type::STR),
                    UnsafeRef::new(Type::dict(
                        UnsafeRef::new(Type::STR),
                        UnsafeRef::new(Type::STR),
                    )),
                ])),
                UnsafeRef::new(Type::dict(
                    UnsafeRef::new(Type::ANY),
                    UnsafeRef::new(Type::ANY),
                )),
            ];
            sup(&types);
        })
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
