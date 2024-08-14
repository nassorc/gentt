module gentt

go 1.22.4

// replace github.com/nassorc/go-codebase/src/sparse_set => ../go-codebase/src/sparse_set

replace github.com/nassorc/go-codebase/src/bitset => ../go-codebase/src/bitset

replace github.com/nassorc/go-codebase/src/ringbuffer => ../go-codebase/src/ringbuffer

require (
	github.com/nassorc/go-codebase/src/bitset v0.0.0-00010101000000-000000000000
	github.com/nassorc/go-codebase/src/ringbuffer v0.0.0-00010101000000-000000000000
// github.com/nassorc/go-codebase/src/sparse_set v0.0.0-00010101000000-000000000000
)
