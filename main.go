package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
}

var db *sql.DB

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, first_name, last_name, username, email, age FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "DB query failed"})
		return
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.Age); err != nil {
			continue
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}
	var u User
	err = db.QueryRow("SELECT id, first_name, last_name, username, email, age FROM users WHERE id=$1", id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.Age)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "DB error"})
		return
	}
	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}
	err := db.QueryRow(
		"INSERT INTO users(first_name, last_name, username, email, age) VALUES($1,$2,$3,$4,$5) RETURNING id",
		u.FirstName, u.LastName, u.Username, u.Email, u.Age,
	).Scan(&u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "insert failed"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}
	res, err := db.Exec(
		"UPDATE users SET first_name=$1, last_name=$2, username=$3, email=$4, age=$5 WHERE id=$6",
		u.FirstName, u.LastName, u.Username, u.Email, u.Age, u.ID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "update failed"})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}
	res, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "delete failed"})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	connStr := os.Getenv("POSTGRES_DSN")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
