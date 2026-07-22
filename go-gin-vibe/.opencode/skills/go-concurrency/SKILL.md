---
name: go-concurrency
description: Use when writing concurrent code, goroutines, channels, or handling parallel operations. Covers goroutine patterns, channel usage, and context propagation.
---

# Go Concurrency Skill

## Rules

1. Always know how goroutine stops
2. Use context for cancellation
3. Use sync.WaitGroup for waiting
4. Use errgroup for concurrent error handling
5. Never start goroutine without knowing how it stops

## Basic Goroutine

```go
// Good - with context for cancellation
func ProcessData(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item // Capture loop variable
        g.Go(func() error {
            return processItem(ctx, item)
        })
    }

    return g.Wait()
}

// Bad - no context, no cancellation
func ProcessData(items []Item) error {
    for _, item := range items {
        go processItem(item) // Leaks if main exits
    }
    return nil
}
```

## Channel Patterns

```go
// Fan-out, fan-in
func ProcessItems(ctx context.Context, items []Item) ([]Result, error) {
    ch := make(chan Result, len(items))
    g, ctx := errgroup.WithContext(ctx)

    // Fan-out
    for _, item := range items {
        item := item
        g.Go(func() error {
            result, err := process(ctx, item)
            if err != nil {
                return err
            }
            ch <- result
            return nil
        })
    }

    // Wait and close channel
    go func() {
        g.Wait()
        close(ch)
    }()

    // Fan-in
    var results []Result
    for result := range ch {
        results = append(results, result)
    }

    return results, g.Err()
}
```

## Buffered Channel as Semaphore

```go
func ProcessWithLimit(ctx context.Context, items []Item, limit int) error {
    sem := make(chan struct{}, limit)
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item
        g.Go(func() error {
            select {
            case sem <- struct{}{}:
                defer func() { <-sem }()
                return process(ctx, item)
            case <-ctx.Done():
                return ctx.Err()
            }
        })
    }

    return g.Wait()
}
```

## Context Propagation

```go
// Always pass context as first parameter
func Handler(c *gin.Context) {
    ctx := c.Request.Context()
    result, err := service.Process(ctx, input)
    // ...
}

func (s *service) Process(ctx context.Context, input Input) (*Output, error) {
    return s.repo.Find(ctx, input.ID)
}

func (r *repository) Find(ctx context.Context, id int64) (*Entity, error) {
    return r.db.WithContext(ctx).First(&entity, id).Error
}
```

## sync.Once

```go
type Singleton struct {
    db *gorm.DB
}

var (
    instance *Singleton
    once     sync.Once
)

func GetSingleton() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            db: connectDB(),
        }
    })
    return instance
}
```

## Best Practices

1. Use errgroup for concurrent operations with error handling
2. Use context for cancellation and timeout
3. Use buffered channels for signaling
4. Never use time.Sleep for synchronization
5. Use sync.WaitGroup for waiting goroutines
6. Use sync.Mutex for shared state
7. Use atomic operations for simple counters
8. Always check for goroutine leaks in tests
