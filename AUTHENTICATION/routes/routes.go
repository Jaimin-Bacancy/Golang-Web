package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"package/controller"
	"package/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var router *mux.Router

//LoadEnvFile is load env files
func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

//CreateRouter is router creation
func CreateRouter() {
	router = mux.NewRouter()
}

//InitializeRoute is add routes
func InitializeRoute() {

	router.HandleFunc("/signup", controller.SingUp).Methods("POST")
	router.HandleFunc("/signin", controller.SignIn).Methods("POST")
	router.HandleFunc("/", controller.Index).Methods("GET")
	router.HandleFunc("/admin", middleware.IsAuthorizedAdmin(controller.AdminIndex)).Methods("GET")
	router.HandleFunc("/admin/display", middleware.IsAuthorizedAdmin(controller.AdminDisplay)).Methods("GET")
	router.HandleFunc("/superadmin", middleware.IsAuthorizedSuperAdmin(controller.SuperAdminIndex)).Methods("GET")
	router.HandleFunc("/user", middleware.IsAuthorizedUser(controller.UserIndex)).Methods("GET")

}

//ServerStart is start server
func ServerStart() {
	serverport := os.Getenv("SERVER_PORT")
	fmt.Println("Server started at http://localhost" + serverport)
	http.ListenAndServe(serverport, router)
}
