package main

import (
	"package/database"
	"package/routes"
)

func main() {
	routes.LoadEnvFile()
	database.Initialmigration()
	routes.CreateRouter()
	routes.InitializeRoute()
	routes.ServerStart()
}
