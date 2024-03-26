package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("PGCONNSTRING"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running on port 3030")
	panic(http.ListenAndServe(":3030", application(NewUsersModel(db))))
}
