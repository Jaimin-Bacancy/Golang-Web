package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"package/controller"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var router *mux.Router

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func CreateRouter() {
	router = mux.NewRouter()
}

func StartServer() {
	serverport := os.Getenv("SERVER_PORT")
	fmt.Println("Server started at http://localhost" + serverport + "/static/index.html")
	http.ListenAndServe(serverport, router)
}

func InitializeRoutes() {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	router.PathPrefix("/static/").Handler(fs)
	http.Handle("/static/", router)
	router.HandleFunc("/calculator", controller.Operation).Methods("POST")
}
