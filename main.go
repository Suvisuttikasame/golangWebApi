package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbName   = "test"
)

var db *sql.DB

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}

	//this allow to access w in order to write file
	fmt.Fprintf(w, "hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Error to read form : message: %v", err)
		return
	}
	fmt.Fprintf(w, "Post request is successful!\n")
	name := r.FormValue("name")
	addr := r.FormValue("address")
	fmt.Fprintf(w, "user name : %s\n", name)
	fmt.Fprintf(w, "user address : %s\n", addr)
	_, err = db.Query(`INSERT INTO mytable(name, address)
	VALUES($1, $2)`, name, addr)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	fmt.Println("Go Server Start here!")
	var err error

	//connect to postgres db
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", host, user, password, port, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connect to postgres DB!")

	//this will look into ./static and always look for index.html
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Starting server at port 8080")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Can not start server at port 8080")
		log.Fatal(err)
	}
}
