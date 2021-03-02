package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"package/database"
	"package/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("secretkeyjwt")

func generatehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	bodydata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("error in body reading")
	}
	var user model.User
	err = json.Unmarshal(bodydata, &user)
	connection := database.GetDatabase()
	defer database.Closedatabase(connection)
	var checkuser model.User
	connection.Where("email = 	?", user.Email).First(&checkuser)
	if checkuser.Email != "" {
		var Error model.Error
		Error.Code = "1001"
		Error.Message = "Email already use"
		bytedata, err := json.MarshalIndent(Error, "", "  ")
		if err != nil {
			log.Fatalln("error in marshal")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytedata)
		return
	}
	user.Password, err = generatehashPassword(user.Password)
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

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

	check := CheckPasswordHash(authdetails.Password, authuser.Password)
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

	validToken, err := GenerateJWT(authuser.Email, authuser.Role)
	var token model.Token
	token.Email = authuser.Email
	token.TokenString = validToken
	if err != nil {
		log.Println("Failed to generate token")
	}
	tokendata, err := json.MarshalIndent(token, "", "  ")
	w.Write([]byte(tokendata))
}

func GenerateJWT(email, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
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
