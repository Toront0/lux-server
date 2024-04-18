package main

import (
	"guthub.com/Toront0/lux-server/internal/api"
	
	
)





func main() {
	
	server := api.NewServer(":3000")

	server.Run()

}