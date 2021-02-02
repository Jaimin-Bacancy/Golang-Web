package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

//Customer is...
type Customer struct {
	Customerid   string
	FirstName    string
	LastName     string
	Email        string
	Dateofbirth  string
	Mobilenumber string
}

var deleteid string
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:1312@localhost/Customer?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully connected to database!")
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		query := "INSERT INTO customer (firstname,lastname,email,dateofbirth,mobilenumber) VALUES($1,$2,$3,$4,$5)"
		customer := Customer{
			FirstName:    r.FormValue("firstname"),
			LastName:     r.FormValue("lastname"),
			Email:        r.FormValue("email"),
			Dateofbirth:  r.FormValue("dateofbirth"),
			Mobilenumber: r.FormValue("mobilenumber"),
		}
		_, err := db.Exec(query, customer.FirstName, customer.LastName, customer.Email, customer.Dateofbirth, customer.Mobilenumber)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/error", 301)
		}
		tmpl.ExecuteTemplate(w, "success.html", struct{ Data string }{"Inserted"})
	}

	http.Redirect(w, r, "/", 301)
}

func getallUser(w http.ResponseWriter, r *http.Request) {

	var customers []Customer
	rows, err := db.Query("SELECT *FROM customer")

	for rows.Next() {
		var customer Customer
		err = rows.Scan(&customer.Customerid, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Dateofbirth, &customer.Mobilenumber)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/error", 301)
		}

		customers = append(customers, customer)
	}

	tmpl.ExecuteTemplate(w, "display.html", customers)
}

func delete(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	deleteid = id
	query := "DELETE FROM customer where id = $1"

	res, err := db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/error", 301)
	}
	count, err := res.RowsAffected()
	if count == 0 {
		http.Redirect(w, r, "/error", 301)
	}
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/error", 301)
	}
	http.Redirect(w, r, "/deletesuccess", 301)
}

func edit(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	query := "SELECT *FROM customer WHERE id=" + id
	rows := db.QueryRow(query)
	var customer Customer
	err := rows.Scan(&customer.Customerid, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Dateofbirth, &customer.Mobilenumber)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/error", 301)
	}

	tmpl.ExecuteTemplate(w, "edit.html", customer)

}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		query := "UPDATE customer SET firstname = $1,lastname = $2,email = $3 ,dateofbirth =$4 ,mobilenumber =$5 where id = $6"
		customer := Customer{
			Customerid:   r.FormValue("id"),
			FirstName:    r.FormValue("firstname"),
			LastName:     r.FormValue("lastname"),
			Email:        r.FormValue("email"),
			Dateofbirth:  r.FormValue("dateofbirth"),
			Mobilenumber: r.FormValue("mobilenumber"),
		}
		_, err := db.Exec(query, customer.FirstName, customer.LastName, customer.Email, customer.Dateofbirth, customer.Mobilenumber, customer.Customerid)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/error", 301)
		}

		tmpl.ExecuteTemplate(w, "success.html", struct{ Data string }{"Updated"})
	}

	http.Redirect(w, r, "/", 301)
}

func deletesuccess(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "deletemsg.html", deleteid)
}

func servererror(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "error.html", nil)
}

func main() {
	defer db.Close()
	fmt.Println("Server started at 7000")
	http.HandleFunc("/", index)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/display", getallUser)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/update", update)
	http.HandleFunc("/error", servererror)
	http.HandleFunc("/deletesuccess", deletesuccess)
	err := http.ListenAndServe(":7000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
