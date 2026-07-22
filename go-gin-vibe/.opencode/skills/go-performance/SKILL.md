---
name: go-performance
description: Use when optimizing Go code, profiling, or improving performance. Covers memory allocation, string operations, and profiling tools.
---

# Go Performance Skill

## Profiling Commands

### CPU Profile

```bash
go test -cpuprofile=cpu.out ./...
go tool pprof cpu.out
```

### Memory Profile

```bash
go test -memprofile=mem.out ./...
go tool pprof mem.out
```

### Benchmark

```bash
go test -bench=. -benchmem ./...
```

### Trace

```bash
go test -trace=trace.out ./...
go tool trace trace.out
```

## Memory Allocation

```go
// Bad - creates new slice each iteration
func merge(slices [][]int) []int {
    var result []int
    for _, s := range slices {
        result = append(result, s...)
    }
    return result
}

// Good - pre-allocate
func merge(slices [][]int) []int {
    size := 0
    for _, s := range slices {
        size += len(s)
    }
    result := make([]int, 0, size)
    for _, s := range slices {
        result = append(result, s...)
    }
    return result
}
```

## String Operations

```go
// Bad - concatenation in loop
func buildString(items []string) string {
    result := ""
    for _, item := range items {
        result += item
    }
    return result
}

// Good - use strings.Builder
func buildString(items []string) string {
    var builder strings.Builder
    builder.Grow(estimatedSize) // Pre-allocate if possible
    for _, item := range items {
        builder.WriteString(item)
    }
    return builder.String()
}
```

## Map Pre-allocation

```go
// Bad
m := make(map[string]int)
for _, item := range items {
    m[item.Key] = item.Value
}

// Good
m := make(map[string]int, len(items))
for _, item := range items {
    m[item.Key] = item.Value
}
```

## Slice Pre-allocation

```go
// Bad
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// Good
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

## Struct Methods

```go
// Bad - copies entire struct
func (u User) GetName() string {
    return u.Name
}

// Good - pointer receiver for large structs
func (u *User) GetName() string {
    return u.Name
}
```

## Interface Assertions

```go
// Bad - type assertion without check
user := entity.(*User)

// Good - type assertion with check
user, ok := entity.(*User)
if !ok {
    return fmt.Errorf("expected *User, got %T", entity)
}
```

## Best Practices

1. Use -benchmem to find allocations
2. Pre-allocate slices and maps when size is known
3. Use strings.Builder for string concatenation
4. Avoid unnecessary allocations
5. Use pointer receivers for large structs
6. Use sync.Pool for frequently allocated objects
7. Profile before optimizing
8. Focus on hot paths first
