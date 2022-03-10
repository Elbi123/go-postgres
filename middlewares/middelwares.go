package middlewares

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"database/sql"
	"encoding/json"
	"net/http"
	"go-postgres/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	 _ "github.com/lib/pq"
 )

 type response struct {
	 ID int64 `json:"id,omitempty"`
	 Message string `json:"message,omitempty"`
 }

 func createConnection() *sql.DB {
	 err := godotenv.Load(".env")

	 if err != nil{
		 log.Fatalf("Error while loading .env file")
	 }

	 db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	 if err != nil {
		 panic(err)
	 }

	 // check the connection

	 err = db.Ping()

	 if err != nil {
		 panic(err)
	 }

	 fmt.Println("Database successfully connected")

	 return db
 }

 func CreateUser(w http.ResponseWriter, r *http.Request) {
	 w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	 w.Header().Set("Access-Control-Allow-Origin", "*")
	 w.Header().Set("Access-Control-Allow-Methods", "POST")
	 w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	 var user models.User

	 // decode the json request to user
	 err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertUser(user)

	res := response{
		ID: insertID,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func insertUser(user models.User) int64 {
	db := createConnection()

	defer db.Close()

	insertUserStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	var id int64

	err := db.QueryRow(insertUserStatement, user.Name, user.Location, user.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int: %v", err)
	}

	user, err := getUser(int64(id))

	if err != nil {
		log.Fatal("Unable to get user: %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func getUser(id int64) (models.User, error) {
	fmt.Println("GET USER is hit\n")
	db := createConnection()

	defer db.Close()

	var user models.User;

	selectStatment := `SELECT * FROM users WHERE userid=$1`

	row := db.QueryRow(selectStatment, id)

	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)

	}

	return user, nil

}


