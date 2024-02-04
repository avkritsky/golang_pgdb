package golang_pgdb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

type User struct {
	ID int
	Username string
}


type Userdata struct {
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


//exists - check user exsists in DB
func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()

	if err != nil {
		fmt.Println("Error in exists:", err)
		return -1
	}
	defer db.Close()

	userID := -1

	statement := fmt.Sprintf(`SELECT "id" FROM "users" WHERE username = '%s'`, username)
	rows, err := db.Query(statement)

	if err != nil {
		fmt.Println("Erorr with execute query:", err)
		return userID
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)

		if err != nil {
			fmt.Println("Error with scan results:", err)
			return userID
		}
		userID = id
	}

	defer rows.Close()

	return userID

}


// AddUser create user,  return ID new user or -1 for error
func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)

	db, err := openConnection()

	if err != nil {
		fmt.Println("Error with open connection:", err)
		return -1
	}

	defer db.Close()

	userID := exists(d.Username)

	if userID != -1 {
		fmt.Println("User already exists!")
		return -1
	}

	insertStatement := `INSERT INTO "users" ("username") values ($1)`

	_, err = db.Exec(insertStatement, d.Username)

	if err != nil {
		fmt.Println("Creating error:", err)
	}

	userID = exists(d.Username)

	if userID == -1 {
		fmt.Println("User not created!")
		return userID
	}

	
}
