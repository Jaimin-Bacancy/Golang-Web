package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

//Image is...
type Image struct {
	ID        string
	Imagepath string
}

var tpl = template.Must(template.ParseGlob("template/*"))

func conn() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:1312@localhost/Image?sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("database connected")
	return db
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func saveimage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, filehandle, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", 301)
		}
		fmt.Println(file)

		imgPath := filepath.Join("static/", filehandle.Filename)
		fmt.Println(imgPath)
		destination, err := os.Create(imgPath)
		if err != nil {
			panic(err)
		}

		defer destination.Close()
		io.Copy(destination, file)
		db := conn()
		defer db.Close()
		query := "INSERT INTO imagetable (imagepath) values($1)"
		_, err = db.Exec(query, imgPath)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/display", 301)
	}
	http.Redirect(w, r, "/", 301)
}

func display(w http.ResponseWriter, r *http.Request) {
	var images []Image
	db := conn()
	query := "SELECT *FROM imagetable"
	rows, err := db.Query(query)
	for rows.Next() {
		var image Image
		err = rows.Scan(&image.ID, &image.Imagepath)
		if err != nil {
			panic(err)
		}
		images = append(images, image)
	}
	tpl.ExecuteTemplate(w, "display.html", images)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/saveimage", saveimage)
	http.HandleFunc("/display", display)
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", fs)
	fmt.Println("server started at 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
