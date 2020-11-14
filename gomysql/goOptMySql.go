package main

import "database/sql"

func main() {
	mysqlUrl := "user:password@tcp(127.0.0.1:3306)/dbname"
	db, err := sql.Open("mysql", mysqlUrl)
	if err != nil {
		panic(err)
	}
	// 注意这行代码要写在上面err判断的下面
	defer db.Close()
}
