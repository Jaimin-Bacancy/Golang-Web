package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"package/controller"
	"package/middleware"

	"github.com/gorilla/handlers"
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
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	router.PathPrefix("/static/").Handler(fs)
	http.Handle("/static/", router)

	router.HandleFunc("/signup", controller.SingUp).Methods("POST")
	router.HandleFunc("/signin", controller.SignIn).Methods("POST")
	router.HandleFunc("/", controller.Index).Methods("GET")
	router.HandleFunc("/admin", middleware.IsAuthorizedAdmin(controller.AdminIndex)).Methods("GET")
	router.HandleFunc("/admin/display", middleware.IsAuthorizedAdmin(controller.AdminDisplay)).Methods("GET")
	router.HandleFunc("/superadmin", middleware.IsAuthorizedSuperAdmin(controller.SuperAdminIndex)).Methods("GET")
	router.HandleFunc("/user", middleware.IsAuthorizedUser(controller.UserIndex)).Methods("GET")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})

}

//ServerStart is start server
func ServerStart() {
	serverport := os.Getenv("SERVER_PORT")
	fmt.Println("Server started at http://localhost" + serverport)
	//http.ListenAndServe(serverport, router)
	http.ListenAndServe(serverport, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
}
