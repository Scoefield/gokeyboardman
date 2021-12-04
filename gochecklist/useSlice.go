package main

import "fmt"

// bad: 未判断data的长度，可导致 index out of range
func badDecode(data []byte) bool {
	if data[0] == 'F' && data[1] == 'U' && data[2] == 'Z' && data[3] == 'Z' && data[4] == 'E' && data[5] == 'R' {
		fmt.Println("Bad")
		return true
	}
	return false
}
	
// bad: slice bounds out of range
func foo() {
	slice := []int{0, 1, 2, 3, 4, 5, 6}
	fmt.Println(slice[:10])
}

// good: 使用data前应判断长度是否合法
func goodDecode(data []byte) bool {
	if len(data) == 6 {
		if data[0] == 'F' && data[1] == 'U' && data[2] == 'Z' && data[3] == 'Z' && data[4] == 'E' && data[5] == 'R' {
			fmt.Println("Good")
			return true
		}
	}
	return false
}

func main() {
	data := []byte{'A', 'B', 'C', 'D'}
	fmt.Println(goodDecode(data))
	fmt.Println(badDecode(data))
}