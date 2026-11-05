Testing checks that your code works correctly. Go has a built-in `testing` package. You do not need to install anything.

Test files end with `_test.go`. To run tests, use:

```bash
go test ./...
```

---

## 1. `*testing.T` — Standard Test

Use `*testing.T` to check that a function returns the correct result. This is for unit tests and integration tests.

Example:

```go
func TestAdd(t *testing.T) {
    t.Parallel() // run this test at the same time as others

    got := Add(2, 3)
    want := 5

    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}
```

Key methods:

- `t.Error()` — mark the test as failed, but keep running
- `t.Errorf()` — mark as failed and print a message
- `t.FailNow()` — fail and stop the test right away
- `t.Skip()` — skip this test
- `t.Parallel()` — run this test at the same time as other parallel tests
- `t.Run()` — run a sub-test

---

## 2. `*testing.B` — Benchmark

Use `*testing.B` to measure how fast a function runs and how much memory it uses.

Run benchmarks with:

```bash
go test -bench=.
```

Example:

```go
func BenchmarkAdd(b *testing.B) {
    b.ReportAllocs() // report memory allocations

    b.ResetTimer() // ignore setup time before the loop
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

Key features:

- `b.N` — Go sets this number automatically (for example 100, 10000, 1000000). Go increases it until the test runs long enough to get a good estimate (about 1 second by default).
- `b.ResetTimer()` — use this if you have heavy setup before the loop, so setup time does not affect the result.
- `b.SetBytes(n)` — measure throughput (MB/s).

---

## 3. `*testing.F` — Fuzz Testing

Fuzz testing was introduced in Go 1.18. It sends random input to your code to find bugs, edge cases, or crashes.

Example:

```go
func FuzzParseJSON(f *testing.F) {
    // 1. Add seed corpus (example inputs)
    f.Add([]byte(`{"key": "value"}`))

    // 2. Run the fuzz target
    f.Fuzz(func(t *testing.T, data []byte) {
        // t here is a normal *testing.T
        _ = ParseJSON(data)
    })
}
```

Key features:

- `f.Add()` — register example input samples
- `f.Fuzz()` — receives a target function. The first argument is a `*testing.T`, followed by the parameter type to test (for example `[]byte`, `string`, `int`)

---

## 4. `testing.TB` — Shared Interface

`testing.TB` is an interface that groups the methods shared by `*T`, `*B`, and `*F` (such as `Log`, `Fatal`, `Error`, `Skip`, `Helper`).

It is useful for making a test helper that works in both normal tests and benchmarks:

```go
// Can be called from Test or Benchmark
func assertNoError(tb testing.TB, err error) {
    tb.Helper() // tell Go to point error lines to the caller of this function
    if err != nil {
        tb.Fatalf("unexpected error: %v", err)
    }
}

func TestFetch(t *testing.T) {
    err := FetchData()
    assertNoError(t, err)
}

func BenchmarkFetch(b *testing.B) {
    err := FetchData()
    assertNoError(b, err)
}
```

---

## Quick Summary

| Type | Main Use | How to Run | What It Does |
|------|----------|------------|--------------|
| `*testing.T` | Unit & integration tests | `go test ./...` | Checks if logic is correct. Pass or fail. |
| `*testing.B` | Benchmarking | `go test -bench=.` | Measures speed and memory usage. |
| `*testing.F` | Fuzz testing | `go test -fuzz=FuzzNama` | Sends random input to find bugs or crashes. |
| `testing.TB` | Shared interface | Internal use | Used in helper functions for tests and benchmarks. |
