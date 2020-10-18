package main

import (
	"fmt"
)

func main() {
	darr := [...]int{57, 89, 90, 82, 100, 78, 67, 69, 59}
	dslice := darr[2:5]
	fmt.Println("array before",darr)
	for i := range dslice {
		dslice[i]++
	}
	fmt.Println("array after",darr)
}