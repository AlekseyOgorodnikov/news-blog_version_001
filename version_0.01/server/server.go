package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	// Closed connect to bd mysql
	defer db.Close()

	/* // Create data in db
	insert, err := db.Query("INSERT INTO `users` (`name`,`age`) VALUES('Bob',25)")
	if err != nil {
		panic(err)
	}
	defer insert.Close() */

	// Read data in mysql server
	res, err := db.Query("SELECT `name`, `age` FROM `users`")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var user User
		err = res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err)
		}
		fmt.Printf("User: %s with age %d.\n", user.Name, user.Age)
	}

	fmt.Println("Connect to mysql server!")
}
