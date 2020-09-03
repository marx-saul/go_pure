package a

//import "fmt"

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
		//fmt.Println("call")		// impure
		//return global_var		// impure
		return global_const	// pure
		//return b.Fb()			// (impure)
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
