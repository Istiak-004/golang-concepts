# What is an Interface?
An interface is a collection of method signatures that define a contract for behavior. Any type that implements all the methods of an interface implicitly satisfies that interface.

## Key Characteristics
1. Implicit Implementation: Types satisfy interfaces automatically
2. Duck Typing: "If it looks like a duck and quacks like a duck, it's a duck"
3. Composable: Interfaces can embed other interfaces
4. First-Class Citizens: Interfaces are types that can be passed around


##  Interface Composition
Interface composition allows you to build new interfaces by combining existing ones through embedding

### Best Practices
1. Prefer small interfaces - Easier to compose and implement
2. Name composed interfaces carefully - Reflect their combined purpose
3. Document expectations - Especially when composed interfaces interact
4. Avoid deep nesting - Typically 2-3 levels is enough
5. Consider interface satisfaction - Ensure types can reasonably implement all methods

## Common Pitfalls
1. Accidental method collisions:

```go
type A interface { Foo() int }
type B interface { Foo() string }
type AB interface { A; B } // Compile error: method Foo conflict
```
2. Over-composition creating interfaces that are too large

3. Unclear requirements when composed interfaces have implicit dependencies


## Interface Segregation Principle in Go
The Interface Segregation Principle (ISP) is a fundamental concept in software design that states that no client should be forced to depend on methods it does not use. In Go, this principle is particularly important due to the language's implicit interface implementation.

### The Problem with Large Interfaces
```go
// Violates ISP
type DocumentProcessor interface {
    Read(file string) ([]byte, error)
    Write(file string, data []byte) error
    Print(data []byte) error
    Fax(data []byte) error
    Scan(file string) ([]byte, error)
}
```

### Segregated Solution
```go
// Properly segregated interfaces
type Reader interface {
    Read(file string) ([]byte, error)
}

type Writer interface {
    Write(file string, data []byte) error
}

type Printer interface {
    Print(data []byte) error
}

type Faxer interface {
    Fax(data []byte) error
}

type Scanner interface {
    Scan(file string) ([]byte, error)
}
```


## Benefits of Interface Segregation in Go
1. Reduced Coupling: Components depend only on what they need
2. Easier Testing: Mock only the required methods
3. Clearer Contracts: Each interface has a single responsibility
4. Better Composition: More flexible to combine small interfaces
5. Easier Refactoring: Changes affect fewer components