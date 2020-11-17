package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string
	Age   int
	Score float32	`json:"-"`
}

// 基本(序列化、反序列化)操作
func testBaseOpt() {
	person := Person{
		"Jack",
		23,
		98.4,
	}
	// struct -> string
	data, err := json.Marshal(&person)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("data: %v \n", string(data))

	// json str -> struct
	var person2 Person
	err = json.Unmarshal(data, &person2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("person2 data: %v", person2)
}

func main() {
	testBaseOpt()
}
