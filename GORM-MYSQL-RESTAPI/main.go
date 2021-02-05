package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

//User is...
type User struct {
	gorm.Model
	Name        string `json:"name"`
	City        string `json:"city"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
}

var muxrouter *mux.Router

func getdatabase() *gorm.DB {
	connection, err := gorm.Open(mysql.Open("root:@/User?parseTime=true"), &gorm.Config{})
	checkerror(err)
	sqldb, err := connection.DB()
	checkerror(err)
	err = sqldb.Ping()
	checkerror(err)
	fmt.Println("connected to database")
	return connection
}

func closedatabase(connection *gorm.DB) {
	sqldb, err := connection.DB()
	checkerror(err)
	sqldb.Close()
}

func initialmigration() {
	connection := getdatabase()
	connection.AutoMigrate(&User{})
	defer closedatabase(connection)
	fmt.Println("migration done")
}

func checkerror(e error) {
	if e != nil {
		panic(e)
	}
}

func getalluser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	connection := getdatabase()
	defer closedatabase(connection)
	var users []User
	connection.Find(&users)
	fmt.Println(users)
	bytedata, err := json.MarshalIndent(users, "", " ")
	checkerror(err)
	w.Write(bytedata)
}

func getoneuser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	connection := getdatabase()
	defer closedatabase(connection)
	var user User
	connection.First(&user, id)
	fmt.Println(user)
	bytedata, err := json.MarshalIndent(user, "", " ")
	checkerror(err)
	w.Write(bytedata)
}

func saveuser(w http.ResponseWriter, r *http.Request) {
	bodydata, err := ioutil.ReadAll(r.Body)
	checkerror(err)
	var user User
	err = json.Unmarshal(bodydata, &user)
	connection := getdatabase()
	defer closedatabase(connection)
	connection.Create(&user)
	bytedata, err := json.MarshalIndent(user, "", "  ")
	checkerror(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}

func updateoneuser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	connection := getdatabase()
	defer closedatabase(connection)
	connection.First(&user, id)
	bodydata, err := ioutil.ReadAll(r.Body)
	checkerror(err)
	err = json.Unmarshal(bodydata, &user)
	checkerror(err)
	connection.Save(&user)
	bytedata, err := json.MarshalIndent(user, "", "  ")
	checkerror(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}

func deleteoneuser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	connection := getdatabase()
	defer closedatabase(connection)
	connection.Delete(&user, id)
	w.Write([]byte("user is deleted"))
}

func createrouter() {
	muxrouter = mux.NewRouter()
}

func initializeroutes() {
	muxrouter.HandleFunc("/user", getalluser).Methods("GET")
	muxrouter.HandleFunc("/user", saveuser).Methods("POST")
	muxrouter.HandleFunc("/user/{id}", getoneuser).Methods("GET")
	muxrouter.HandleFunc("/user/{id}", updateoneuser).Methods("PUT")
	muxrouter.HandleFunc("/user/{id}", deleteoneuser).Methods("DELETE")
}

func serverstart() {
	fmt.Println("Server started at 9123")
	http.ListenAndServe(":9123", muxrouter)
}

func main() {
	initialmigration()
	createrouter()
	initializeroutes()
	serverstart()
}
