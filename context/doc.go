// The context is created to hold all the required data from previous steps to be reused.
// Its goal preventing from holding a global state if it's not needed what improves maintainability.
// The context has a set of helper functions which simplify the end-user code.
//  ctx.Set(name{}, "John")
//  fmt.Printf("Hi %s\n", ctx.GetString(name{})) // prints "Hi John"
//
// Good practices
//
// It’s a good practice to use custom structs as keys instead of strings or any built-in types
// to avoid collisions between steps using context.

package context
