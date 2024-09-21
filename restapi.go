package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	log.Println(r.Body)
	fmt.Fprintf(w, "Hello World!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login")
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
