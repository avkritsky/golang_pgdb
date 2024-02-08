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

	insertStatement = `insert into "userdata" ("user_id", "name", "surname", "description") values ($1, $2, $3, $4)`

	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)

	if err != nil {
		fmt.Println("Error with insert #2:", err)
		return -1
	}

	return userID

}

func DeleteUser(id int) error {
	db, err := openConnection()

	if err != nil {
		fmt.Println("Connection error:", err)
		return err
	}

	defer db.Close()

	checkQuery := `SELECT "username" FROM users WHERE "id" = $1`

	rows, err := db.Query(checkQuery, id)

	if err != nil {
		fmt.Println("Error check:", err)
		return err
	}

	var username string

	for rows.Next() {
		err = rows.Scan(&username)

		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with id #%d does not exist!\n", id)
	}

	deleteStatement := `DELETE FROM "users" WHERE "id" = $1`
	_, err = db.Exec(deleteStatement, id)

	if err != nil {
		return err
	}

	deleteStatement = `DELETE FROM "userdata" WHERE "userid" = $1`
	_, err = db.Exec(deleteStatement, id)

	if err != nil {
		return err
	}

	return nil
}

func ListUser() ([]Userdata, error) {
	db, err := openConnection()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	data := []Userdata{}

	rows, err := db.Query(`SELECT "id", username, name, surname, description  FROM "users", "userdata" WHERE users.id = userdata.userid`)

	for rows.Next() {
		var id int
		var (
			username string
	        name string
			surname string
			description string
		)

		err = rows.Scan(&id, &username, &name, &surname, &description)

		if err != nil {
			return data, err
		}

		temp := Userdata{
			id, username, name, surname, description,
		}

		data = append(data, temp)
	}

	defer rows.Close()

	return data, nil
}
