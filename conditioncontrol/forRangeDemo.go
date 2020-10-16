package main

import "fmt"

func main() {
	sliceData := []int{1, 2, 3, 4, 5}
	for index, val := range sliceData {
		fmt.Println("index=", index, "value=", val)
	}
}
