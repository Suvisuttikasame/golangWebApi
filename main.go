package main

import (
	"database/sql"
	"encoding/json"
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

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Error to read form : message: %v", err)
		return
	}
	fmt.Fprintf(w, "Post request is successful!\n")
	topic := r.FormValue("topic")
	note := r.FormValue("note")
	rating := r.FormValue("emotion")
	fmt.Fprintf(w, "user topic : %s\n", topic)
	fmt.Fprintf(w, "user note : %s\n", note)
	fmt.Fprintf(w, "user rating : %s\n", rating)
	_, err = db.Query(`INSERT INTO mytable(name, address, rating)
	VALUES($1, $2, $3)`, topic, note, rating)
	if err != nil {
		log.Fatal(err)
	}

}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Id      int
		Name    string
		Address string
		Rating  int64
	}
	var datas []data
	var (
		id      int
		name    sql.NullString
		address sql.NullString
		rating  sql.NullInt64
	)
	fmt.Println("Starting to get all data")
	rows, err := db.Query(`SELECT * from mytable`)
	defer rows.Close()
	if err != nil {
		fmt.Println("Fail to get all data")
		log.Fatal(err)
		return
	}
	for rows.Next() {
		err := rows.Scan(&id, &name, &address, &rating)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name.String)
		datas = append(datas, data{
			Id:      id,
			Name:    name.String,
			Address: address.String,
			Rating:  rating.Int64,
		})
	}
	js, err := json.Marshal(datas)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(js))
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
	http.HandleFunc("/get-all-notes", getAllHandler)

	fmt.Println("Starting server at port 8080")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Can not start server at port 8080")
		log.Fatal(err)
	}
}
