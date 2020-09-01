package a

//import "fmt"
//import "b"

// Hello!
// I am a.f()
func f() {
	var gopher int32
	print(gopher)
}

var global_var int32 = 1
const global_const = 1

// @pure
func factorial(n int32) int32 {
	if n <= 0 {
		return global_var		// impure
		//return global_const	// pure
		//return b.Fb()			// (impure)
	} else {
		return n * factorial(n-1)
	}
}

