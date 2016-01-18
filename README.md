## TtlMap - an in-memory TTL map with expiration

1. Thread-safe
2. Auto-Expiring after a certain time
3. Auto-Extending expiration on `Get`s

#### Usage
```go
import (
  "time"
  "github.com/codinl/ttlmap"
)

func main () {
  m := ttlmap.NewTtlMap(time.Second)
  m.Set("key", "value")
  value, exists := m.Get("key")
  count := m.Count()
}
```