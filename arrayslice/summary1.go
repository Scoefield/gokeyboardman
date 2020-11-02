package main

import "fmt"

func sliceFunc(s []int) {
	for i := range s {
		s[i] += 1
	}
}

func main() {
	s := []int{1, 3, 5}
	sliceFunc(s)
	fmt.Println(s)
}
