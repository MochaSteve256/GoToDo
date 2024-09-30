package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	mathrand "math/rand"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	secretKey = "your-32-byte-secret-key-here-123"
)

// Functions

func generateSessionToken() string {
	b := make([]byte, 255)
	for i := range b {
		b[i] = byte(mathrand.Intn(62) + 48)
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

func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", err
	}
	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
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

		// Encrypt email
		encryptedEmail, err := encrypt(email)
		if err != nil {
			log.Printf("Error encrypting email: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Hash email for lookup purposes
		emailHash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing email: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Insert user into userDB
		_, err = db.Exec("INSERT INTO users (email, email_hash, password_hash, user_name) VALUES ($1, $2, $3, $4);",
			encryptedEmail, emailHash, passwordHash, username)
		if err != nil {
			log.Printf("Error inserting user into database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
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
	r.HandleFunc("/login", login(db))
	r.HandleFunc("/getAllUsers", getUsers(db))
	r.HandleFunc("/register", registerUser(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	// Connect to database (on windows with local ip, on mac routed through tailscale vpn
	os := runtime.GOOS
	var ipStr string
	switch os {
	case "windows":
		ipStr = "192.168.178.73"
	case "darwin":
		ipStr = "100.96.62.40"
	}
	connStr := "postgres://gotodoapp:znoi4WfR6@" + ipStr + ":5432/gotodo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create http server
	handleRequests(db)

}
