# gox

A collection of Go utilities providing AWS-style identifiers and slice operations with generics support.

## Packages

### `idx` - AWS-Style Identifiers

Generate AWS-style identifiers with flexible namespace-based object IDs in the format: `environment:type:object_id`

#### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/letmevibethatforyou/gox/idx"
)

// Declare types at package level
const (
    UserType     idx.Type = "user"
    ProductType  idx.Type = "product"
    OrderType    idx.Type = "order"
    SessionType  idx.Type = "session"
)

func main() {
    // Create namespace (prd and empty string become "vibe")
    ns := idx.NewNamespace("prd")
    
    // Generate IDs with auto-generated values
    userID, err := ns.NewID(UserType)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(userID.String()) // "vibe:user:auto_generated_ksuid"
    
    // Create IDs with custom values
    customID, err := ns.NewIDWithValue(ProductType, "custom123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(customID.String()) // "vibe:product:custom123"
    
    // Access ID components
    fmt.Printf("Env: %s, Type: %s, Value: %s\n", 
        userID.Env(), userID.Type(), userID.Value())
}
```

#### Advanced Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/letmevibethatforyou/gox/idx"
)

// Define all types at package level for reuse
const (
    AdminUserType   idx.Type = "admin_user"
    RegularUserType idx.Type = "regular_user"
    APIKeyType      idx.Type = "api_key"
    WebhookType     idx.Type = "webhook"
)

func main() {
    // Different environments
    devNS := idx.NewNamespace("dev")
    stagingNS := idx.NewNamespace("staging")
    prodNS := idx.NewNamespace("prd") // becomes "vibe"
    
    // Create different types of IDs
    adminID, _ := devNS.NewID(AdminUserType)
    apiKeyID, _ := stagingNS.NewIDWithValue(APIKeyType, "key_12345")
    webhookID, _ := prodNS.NewID(WebhookType)
    
    fmt.Printf("Admin: %s\n", adminID.String())      // "dev:admin_user:..."
    fmt.Printf("API Key: %s\n", apiKeyID.String())   // "staging:api_key:key_12345"
    fmt.Printf("Webhook: %s\n", webhookID.String())  // "vibe:webhook:..."
    
    // Parse existing IDs
    parsed, err := idx.ParseID("vibe:user:abc123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Parsed - Env: %s, Type: %s, Value: %s\n",
        parsed.Env(), parsed.Type(), parsed.Value())
        
    // Type validation
    validType, err := idx.ParseType("valid_type_123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Valid type: %s\n", validType)
    
    // Invalid type example
    if _, err := idx.ParseType("123invalid"); err != nil {
        fmt.Printf("Invalid type error: %v\n", err)
    }
}
```

### `slicex` - Generic Slice Operations

Utility functions for common slice operations with full generics support.

#### Basic Operations

```go
package main

import (
    "fmt"
    
    "github.com/letmevibethatforyou/gox/slicex"
)

// Define types at package level
type User struct {
    ID   int
    Name string
    Age  int
}

const (
    MinAdultAge = 18
)

var (
    users = []User{
        {ID: 1, Name: "Alice", Age: 30},
        {ID: 2, Name: "Bob", Age: 25},
        {ID: 3, Name: "Alice", Age: 30}, // duplicate
        {ID: 4, Name: "Charlie", Age: 16},
        {ID: 0, Name: "", Age: 0}, // zero value
    }
)

func main() {
    // Remove duplicates
    names := []string{"Alice", "Bob", "Alice", "Charlie"}
    uniqueNames := slicex.Unique(names)
    fmt.Printf("Unique names: %v\n", uniqueNames) // [Alice Bob Charlie]
    
    // Filter non-zero values
    numbers := []int{0, 1, 2, 0, 3, 0, 4}
    nonZero := slicex.FilterNonZero(numbers)
    fmt.Printf("Non-zero: %v\n", nonZero) // [1 2 3 4]
    
    // Map transformation
    ages := slicex.Map(users, func(u User) int { return u.Age })
    fmt.Printf("Ages: %v\n", ages) // [30 25 30 16 0]
    
    // Group by criteria
    grouped := slicex.Group(users, func(u User) string { return u.Name })
    fmt.Printf("Grouped by name: %v\n", grouped)
}
```

#### Concurrent Processing

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/letmevibethatforyou/gox/slicex"
)

// Define processing function at package level
var processUser = func(ctx context.Context, userID int) (string, error) {
    // Simulate some work
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("processed_user_%d", userID), nil
}

const (
    MaxConcurrency = 3
    ProcessTimeout = 5 * time.Second
)

var userIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), ProcessTimeout)
    defer cancel()
    
    // Configure and execute concurrent processing
    results, err := slicex.MapConcurrent(processUser).
        WithConcurrency(MaxConcurrency).
        WithStopOnError(false).
        Execute(ctx, userIDs)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Processed %d users concurrently\n", len(results))
    for i, result := range results {
        fmt.Printf("User %d -> %s\n", userIDs[i], result)
    }
}
```

## Installation

```bash
go get github.com/letmevibethatforyou/gox
```

## Requirements

- Go 1.24+
- Dependencies:
  - `github.com/segmentio/ksuid` (for unique ID generation)

## Features

### `idx` Package Features
- ✅ AWS-style identifier format (`environment:type:object_id`)
- ✅ Namespace-based ID generation
- ✅ Type validation and parsing
- ✅ Custom and auto-generated values
- ✅ Environment normalization (prd → vibe)

### `slicex` Package Features
- ✅ Generic slice operations (Go 1.18+)
- ✅ Duplicate removal with order preservation
- ✅ Zero-value filtering
- ✅ Map transformations
- ✅ Grouping operations
- ✅ Concurrent processing with configurable workers
- ✅ Error handling and context cancellation support

## License

MIT