package main

import (
	"package/routes"
)

func main() {
	routes.LoadEnvFile()
	routes.CreateRouter()
	routes.InitializeRoutes()
	routes.StartServer()
}
