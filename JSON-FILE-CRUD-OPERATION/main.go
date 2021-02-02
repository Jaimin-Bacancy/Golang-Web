package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	uuid "github.com/nu7hatch/gouuid"
)

type User struct {
	Userid       string
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
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	}
}

func insert(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		usersfiledata, err := ioutil.ReadFile("user.json")
		if err != nil {
			fmt.Println(err)
		}
		Map := make(map[string]bool)

		var alluserdetails []User
		err = json.Unmarshal([]byte(usersfiledata), &alluserdetails)
		if err != nil {
			fmt.Println("Error JSON Unmarshling for user file")
			fmt.Println(err)

		}

		for _, user := range alluserdetails {
			Map[user.Userid] = true
		}

		userid, err := uuid.NewV4()
		fmt.Println(userid)

		userformdetails := User{
			Userid:       userid.String(),
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
		t.Execute(w, struct {
			Data string
		}{"Inserted"})
	}
}

func deleteuserdata(w http.ResponseWriter, r *http.Request) {

	userid := r.URL.Query().Get("id")
	var alluserdata []User

	allfiledata, err := ioutil.ReadFile("user.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(allfiledata, &alluserdata)
	if err != nil {
		fmt.Println(err)
	}

	finaldata := make([]User, 0)

	for _, user := range alluserdata {
		if user.Userid != userid {
			finaldata = append(finaldata, user)
		}
	}

	finaljsondata, err := json.MarshalIndent(finaldata, "", "")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("user.json", finaljsondata, 0644)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/display", 301)

}

func editform(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("edit.html")
	if err != nil {
		fmt.Println(err)
	}

	userid := r.URL.Query().Get("id")
	allfiledata, err := ioutil.ReadFile("user.json")

	if err != nil {
		fmt.Println(err)
	}
	var users []User
	err = json.Unmarshal(allfiledata, &users)
	if err != nil {
		fmt.Println(err)
	}

	var updatedata User
	for _, user := range users {
		if user.Userid == userid {
			updatedata = user
			break
		}
	}
	t.Execute(w, updatedata)
}

func edituserdata(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userformdetails := User{
			Userid:       r.FormValue("id"),
			Firstname:    r.FormValue("firstname"),
			Lastname:     r.FormValue("lastname"),
			DataofBirth:  r.FormValue("dateofbirth"),
			Email:        r.FormValue("email"),
			MobileNumber: r.FormValue("mobilenumber"),
		}

		filedata, err := ioutil.ReadFile("user.json")
		if err != nil {
			fmt.Println(err)
		}

		var users []User

		err = json.Unmarshal(filedata, &users)
		if err != nil {
			fmt.Println(err)
		}

		var finaluserdetails []User
		for _, user := range users {
			if user.Userid == r.FormValue("id") {
				user = userformdetails
				fmt.Println(user)
				fmt.Println(userformdetails)
			}
			finaluserdetails = append(finaluserdetails, user)
		}

		finaluserfiledetails, err := json.MarshalIndent(finaluserdetails, "", "")
		err = ioutil.WriteFile("user.json", finaluserfiledetails, 0644)
		if err != nil {
			fmt.Println(err)
		}

		t, err := template.ParseFiles("success.html")
		if err != nil {
			fmt.Println(err)
		}

		t.Execute(w, struct {
			Data string
		}{"Updated"})
	}

}

func main() {
	fmt.Println("Server started at 8080")
	http.HandleFunc("/", index)
	http.HandleFunc("/dataprocess", insert)
	http.HandleFunc("/display", getalldata)
	http.HandleFunc("/edit", editform)
	http.HandleFunc("/update", edituserdata)
	http.HandleFunc("/delete", deleteuserdata)
	http.ListenAndServe(":8080", nil)

}
