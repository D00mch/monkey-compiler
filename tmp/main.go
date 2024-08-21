package main

import (
	"dumch/monkey/object"
	"fmt"
)

type Bob struct {
	name string
}

func main() {
	b1 := Bob {name: "Bob"}
	b2 := Bob {name: "Bob"}
	fmt.Printf("%v", b1 == b2)

	var result object.ReturnValue
	fmt.Printf("%v", result.Type())
}
