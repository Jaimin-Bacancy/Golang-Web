package main

import (
	"package/database"
	"package/routes"
)

func main() {
	database.Initialmigration()
	routes.LoadEnvFile()
	routes.CreateRouter()
	routes.InitializeRoute()
	routes.ServerStart()
}
