package main

import (
	"fmt"
)

func main() {
        str:= HashCode("NSPDM_HSQL_R01_TEMP_S09")
	fmt.Println(str)
}
func HashCode(key string) int {
	var index int = 0
	index = int(key[0])
	for k := 0; k < len(key); k++ {
		index *= (1103515245 + int(key[k]))
	}
        fmt.Println(index)
	index >>= 27
        fmt.Println(index)
	index &= 5 - 1
	return index
}
