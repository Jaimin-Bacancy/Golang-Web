package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type User struct {
	Userid       int
	Firstname    string
	Lastname     string
	DataofBirth  string
	Email        string
	MobileNumber string
}

func getalldata(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		usersfiledata, err := ioutil.ReadFile("user.json")
		if err != nil {
			fmt.Println(err)
		}

		var alluserdetails []User
		err = json.Unmarshal([]byte(usersfiledata), &alluserdetails)
		if err != nil {
			fmt.Println("Error JSON Unmarshling for user file")
			fmt.Println(err)

		}
		t, _ := template.ParseFiles("display.html")
		t.Execute(w, alluserdetails)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("form.html")
		t.Execute(w, nil)
	}
}

func formhandling(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		usersfiledata, err := ioutil.ReadFile("user.json")
		if err != nil {
			fmt.Println(err)
		}
		Map := make(map[int]bool)

		var alluserdetails []User
		err = json.Unmarshal([]byte(usersfiledata), &alluserdetails)
		if err != nil {
			fmt.Println("Error JSON Unmarshling for user file")
			fmt.Println(err)

		}

		for _, user := range alluserdetails {
			Map[user.Userid] = true
		}
		i := 1
		userid := 0
		for {
			if _, ok := Map[i]; !ok {
				userid = i
				break
			}
			i++
		}

		userformdetails := User{
			Userid:       userid,
			Firstname:    r.FormValue("firstname"),
			Lastname:     r.FormValue("lastname"),
			DataofBirth:  r.FormValue("dateofbirth"),
			Email:        r.FormValue("email"),
			MobileNumber: r.FormValue("mobilenumber"),
		}

		alluserdetails = append(alluserdetails, userformdetails)
		file, _ := json.MarshalIndent(alluserdetails, "", " ")
		_ = ioutil.WriteFile("user.json", file, 0644)
		t, _ := template.ParseFiles("success.html")
		t.Execute(w, nil)

	}
}

func main() {
	fmt.Println("User Form Handling")
	http.HandleFunc("/", index)
	http.HandleFunc("/dataprocess", formhandling)
	http.HandleFunc("/display", getalldata)
	fmt.Println("Server started at 8080")
	http.ListenAndServe(":8080", nil)

}
