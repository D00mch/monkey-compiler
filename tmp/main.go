package main

import (
	"dumch/monkey/code"
	"encoding/binary"
	"fmt"
)

func main() {

	instruction := encode(1)
	instruction1 := encode(256)
	instruction2 := encode(512)
	instruction3 := encode(65534)

	fmt.Printf("B: %v\n", instruction)
	fmt.Printf("B: %v\n", instruction1)
	fmt.Printf("B: %v\n", instruction2)
	fmt.Printf("B: %v\n", instruction3)

	println(code.OpConstant)
}

func encode(n int) []byte {
	instruction := make([]byte, 4)
	binary.BigEndian.PutUint16(instruction, uint16(n))
	return instruction
}
