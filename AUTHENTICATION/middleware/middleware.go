package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func IsAuthorizedAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		secretkey := os.Getenv("SECRET_KEY")
		var mySigningKey = []byte(secretkey)
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["role"])
			if claims["role"] == "admin" {
				handler.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

	}
}

func IsAuthorizedSuperAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var mySigningKey = []byte("secretkeyjwt")
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["role"])
			if claims["role"] == "superadmin" {
				handler.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

	}
}

func IsAuthorizedUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var mySigningKey = []byte("secretkeyjwt")
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["role"])
			if claims["role"] == "user" {
				handler.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

	}
}
