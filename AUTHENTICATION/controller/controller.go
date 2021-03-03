package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"package/database"
	"package/model"
	"package/utility"
)

func SingUp(w http.ResponseWriter, r *http.Request) {

	bodydata, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln("error in body reading")
	}

	var user model.User
	err = json.Unmarshal(bodydata, &user)
	if err != nil {
		log.Fatalln("error in unmarshal")
	}

	connection := database.GetDatabase()
	defer database.Closedatabase(connection)

	var checkuser model.User
	connection.Where("email = 	?", user.Email).First(&checkuser)

	//check email is alredy register
	if checkuser.Email != "" {
		var Error model.Error
		Error.Code = "1001"
		Error.Message = "Email already use"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error)
		return
	}

	user.Password, err = utility.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}
	connection.Create(&user)
	bytedata, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatalln("error in marshal")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bodydata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("error in body reading")
	}
	var authdetails model.Authentatication
	err = json.Unmarshal(bodydata, &authdetails)
	connection := database.GetDatabase()
	defer database.Closedatabase(connection)
	var authuser model.User
	connection.Where("email = 	?", authdetails.Email).First(&authuser)

	if authuser.Email == "" {
		var Error model.Error
		Error.Code = "1002"
		Error.Message = "Username or Passwrod is wrong"
		errordata, err := json.MarshalIndent(Error, "", "  ")
		if err != nil {
			log.Fatalln("error in marshal")
		}
		w.Write(errordata)
		return
	}

	check := utility.CheckPasswordHash(authdetails.Password, authuser.Password)
	if !check {
		var Error model.Error
		Error.Code = "1002"
		Error.Message = "Username or Passwrod is wrong"
		errordata, err := json.MarshalIndent(Error, "", "  ")
		if err != nil {
			log.Fatalln("error in marshal")
		}
		w.Write(errordata)
		return
	}

	validToken, err := utility.GenerateJWT(authuser.Email, authuser.Role)
	var token model.Token
	token.Email = authuser.Email
	token.TokenString = validToken
	if err != nil {
		log.Println("Failed to generate token")
	}
	tokendata, err := json.MarshalIndent(token, "", "  ")
	w.Write([]byte(tokendata))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ADMIN INDEX PAGE"))
}

func SuperAdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SuperAdminIndex INDEX PAGE"))
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User INDEX PAGE"))
}
