package main

import (
	"postgre/routes"
)

func main() {
	routes.Routes()
	routes.StartServer()
}
