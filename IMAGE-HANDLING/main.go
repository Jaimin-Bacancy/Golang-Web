package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var tpl = template.Must(template.ParseGlob("template/*"))

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func saveimage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, filehandle, err := r.FormFile("file")
		fmt.Println(file)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		imgPath := filepath.Join("profiles/", filehandle.Filename)
		fmt.Println(imgPath)
		destination, err := os.Create(imgPath)
		if err != nil {
			panic(err)
		}

		defer destination.Close()
		io.Copy(destination, file)
		fmt.Fprint(w, "image is save")
	}
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/saveimage", saveimage)
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./profiles")))
	http.Handle("/static/", fs)
	fmt.Println("server started at 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
