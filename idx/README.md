# Identity Package

The `identity` package provides AWS-style identifiers with flexible object ID generation. It implements a namespace-based ID system where IDs follow the format: `environment:type:object_id`.

## Features

- **AWS-style IDs**: Structured identifiers with environment, type, and object ID components
- **Flexible Object IDs**: Auto-generate unique values or provide your own custom values
- **Namespace-based API**: Environment-scoped ID creation through namespaces
- **Special "Vibe" Logic**: "prd" and empty environments are automatically normalized to "vibe"
- **Type Safety**: Strong typing for object types with validation
- **Full Roundtrip**: Parse IDs from strings and serialize back to strings

## Installation

```bash
go get github.com/letmevibethatforyou/gox/idx
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/letmevibethatforyou/gox/identity"
)

func main() {
    // Create a namespace
    ns := identity.NewNamespace("dev")
    
    // Create an ID with auto-generated value
    userType := identity.Type("user")
    id, err := ns.NewID(userType)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(id.String()) // "dev:user:auto_generated_unique_value"
    
    // Create an ID with custom value
    customID, err := ns.NewIDWithValue(userType, "custom123")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(customID.String()) // "dev:user:custom123"
    
    // Parse an existing ID
    parsed, err := identity.ParseID("dev:user:custom123")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Environment:", parsed.Env())   // "dev"
    fmt.Println("Type:", parsed.Type())         // "user"
    fmt.Println("Value:", parsed.Value())       // "custom123"
}
```

### Special "Vibe" Environment

The package includes special handling for production environments:

```go
// These all create "vibe" namespace
prodNS := identity.NewNamespace("prd")  // becomes "vibe"
emptyNS := identity.NewNamespace("")    // becomes "vibe"

id, _ := prodNS.NewID(identity.Type("user"))
fmt.Println(id.String()) // "vibe:user:auto_generated_value"
```

### Custom Types

You can define custom object types that meet the validation requirements:

```go
// Valid custom types
customType := identity.Type("order_item")     // Valid: letters, numbers, underscores
apiKeyType := identity.Type("API_Key")        // Valid: starts with letter
sessionType := identity.Type("session123")   // Valid: alphanumeric

// Create ID with custom type
ns := identity.NewNamespace("staging")
id, err := ns.NewID(customType)
```

## API Reference

### Types

#### `Type`
Represents an object type identifier. Must follow these rules:
- Start with a letter (a-z, A-Z)
- Contain only letters, numbers, and underscores
- Maximum 32 characters
- Cannot contain colons
- Cannot be empty

#### `Namespace`
Represents an environment context for creating IDs.

#### `ID`
Represents a complete identifier with environment, type, and object ID components.

### Functions

#### `NewNamespace(environment string) Namespace`
Creates a new namespace with special handling for "prd" and empty environments (both become "vibe").

#### `ParseType(s string) (Type, error)`
Creates and validates a Type from a string.

#### `ParseID(s string) (ID, error)`
Parses a string representation of an ID in the format `environment:type:object_id`.

### Methods

#### `Namespace.NewID(objectType Type) (ID, error)`
Creates a new ID within the namespace using the specified object type and an auto-generated unique value.

#### `Namespace.NewIDWithValue(objectType Type, value string) (ID, error)`
Creates a new ID within the namespace using the specified object type and a custom value provided by the caller.

#### `Namespace.Environment() string`
Returns the normalized environment name for the namespace.

#### `ID.Env() string`
Returns the environment component of the ID.

#### `ID.Type() Type`
Returns the object type component of the ID.

#### `ID.Value() string`
Returns the object ID component of the ID.

#### `ID.String() string`
Returns the full string representation in the format `environment:type:object_id`.

#### `ID.Validate() error`
Validates that all components of the ID are valid.

#### `Type.String() string`
Returns the string representation of the Type.

#### `Type.Validate() error`
Validates that the Type meets all requirements.

## Examples

### Multiple Environments

```go
// Different environments
devNS := identity.NewNamespace("dev")
stagingNS := identity.NewNamespace("staging")
prodNS := identity.NewNamespace("prd") // becomes "vibe"

userType := identity.Type("user")

devID, _ := devNS.NewID(userType)
fmt.Println(devID.String()) // "dev:user:..."

stagingID, _ := stagingNS.NewID(userType)
fmt.Println(stagingID.String()) // "staging:user:..."

prodID, _ := prodNS.NewID(userType)
fmt.Println(prodID.String()) // "vibe:user:..."
```

### Error Handling

```go
// Invalid type
invalidType := identity.Type("123invalid") // starts with number
_, err := identity.ParseType("123invalid")
if err != nil {
    fmt.Println("Type validation failed:", err)
}

// Invalid ID format
_, err = identity.ParseID("invalid:format")
if err != nil {
    fmt.Println("ID parsing failed:", err)
}
```

### Custom Values

```go
// Use your own values
ns := identity.NewNamespace("prd") // becomes "vibe"
customID, _ := ns.NewIDWithValue(identity.Type("order"), "ord_12345")
fmt.Println(customID.String()) // "vibe:order:ord_12345"

// Auto-generated values
autoID, _ := ns.NewID(identity.Type("session"))
fmt.Println(autoID.String()) // "vibe:session:auto_generated_unique_value"
```

## Dependencies

This package uses [github.com/segmentio/ksuid](https://github.com/segmentio/ksuid) internally for auto-generating unique values, but this is an implementation detail that users don't need to be concerned with.