package b

import "fmt"

const constant = 1

// @pure
func Fb() int32 {
	fmt.Println("f0 called");
	return constant
}

