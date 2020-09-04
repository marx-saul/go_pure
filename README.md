# go_pure
A static analysis tool for go-lang that ensures designated functions are pure.

Functions decorated by <code>@pure</code> is throw to purity check.

The norm of function purity is :

* Do not refer to mutable variables that is declared outside the function scope
* Do not call impure functions

These implies that the function always returns the same result when same arguments are passed.

Note that pure functions **can have side effect**, it can change the content of arguments of pointer type.
It also allows <code>print</code>.

## Example
```go
package a

import "fmt"

// @pure
func f() {
	var gopher int32
	print(gopher)
}

var global_var int32 = 1
const global_const = 1

// @pure
func factorial(n int32) int32 {
	if n <= 0 {
		//fmt.Println("call")	// impure
		//return global_var		// impure
		return global_const		// pure
	} else {
		return n * factorial(n-1)
	}
}


// @pure
func fibonacci(n int) int {
	if n <= 0 {
		return 0;
	}
	
	previous := 0
	current  := 1
	
	for i := 1 ; i <= n ; i++ {
		next := previous + current
		previous = current
		current = next
	}
	return current
}
```

## Install
```sh
$ go get github.com/marx-saul/go_pure/cmd/go_pure
```

## Usage
```sh
$ go vet -vettool=`which go_pure` your-package
```

