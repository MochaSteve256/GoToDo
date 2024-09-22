package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Functions

func generateSessionToken() string {
	b := make([]byte, 255)
	for i := range b {
		b[i] = byte(rand.Intn(62) + 48)
		if b[i] > 57 && b[i] < 65 {
			b[i] += 39
		}
	}
	return string(b)
}

func getUserIDfromToken(token string, db *sql.DB) int {
	// Query tokenDB for validating token and gathering userID

	return 0
}

// Endpoints
func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	log.Println(r.Body)
	fmt.Fprintf(w, "Hello World!")
}

func registerUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get email, username, and password from request
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		//encrypt email

		//hash email

		//hash password

		// Insert user into userDB
		_, err := db.Exec("INSERT INTO users (email,
    										email_hash,
										    password_hash,
										    user_name
				) VALUES ();", email, username, password)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get username and password from request
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Query userDB for validating username and passwordHash
		rows, err := db.Query("SELECT * FROM users WHERE username = $1 AND passwordHash = $2", username, password)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()
		// Loop through rows and process the data
		for rows.Next() {
			var id int
			var name string
			var email string
			var password string
			if err := rows.Scan(&id, &name, &email, &password); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "{\"token\": \"%s\"}", generateSessionToken())
		}

	}
}

func getTodos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Query all todos
		rows, err := db.Query("SELECT * FROM todos")
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()
		// Loop through rows and process the data
		for rows.Next() {
			var id int
			var title string
			var completed bool
			if err := rows.Scan(&id, &title, &completed); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "ID: %d, Title: %s, Completed: %t\n", id, title, completed)
		}
	}
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query all users
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()
		// Loop through rows and process the data
		for rows.Next() {
			var id int
			var name string
			var email string
			var password string
			if err := rows.Scan(&id, &name, &email, &password); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "ID: %d, Name: %s, Email: %s\n", id, name, email)
		}
	}
}

func handleRequests(db *sql.DB) {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/getAllUsers", getUsers(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	// Connect to database
	connStr := "postgres://gotodoapp:znoi4WfR6@192.168.178.73:5432/gotodo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create http server
	handleRequests(db)

}
