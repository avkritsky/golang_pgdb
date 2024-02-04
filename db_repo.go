package golang_pgdb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	ID int
	Username string
}


type Userdate struct {
	ID int
	Username string
	Name string
	Surname string
	Description string
}

var (
	Hostname = ""
	Port = 2345
	Login = ""
	Password = ""
	Database = ""
)


func openConnection() (*sql.DB, error) {
	// connection str params
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Login, Password, Database)
	
	// open connection
	db, err := sql.Open("postgres", conn)
	
	// check error
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil, err
	}
	
	return db, nil
}