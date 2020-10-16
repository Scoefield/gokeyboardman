package main

import "fmt"

func main() {
	// 定义局部变量，class等级，marks分数
	var class string = "B"
	var marks int = 90

	switch marks {
		case 90: class = "A"
		case 80: class = "B"
		case 50,60,70 : class = "C"  // case 后可以由多个数值
		default: class = "D"		// 上面case都匹配不到时，会执行default的代码
	}

	switch {
		case class == "A" :
			fmt.Printf("优秀!\n" )
		case class == "B", class == "C" :
			fmt.Printf("良好\n" )
		case class == "D" :
			fmt.Printf("及格\n" )
		case class == "F":
			fmt.Printf("不及格\n" )
		default:
			fmt.Printf("差\n" )
	}
	fmt.Printf("你的等级是 %s\n", class)
}