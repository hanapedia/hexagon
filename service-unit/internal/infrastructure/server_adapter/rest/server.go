package rest

import (
	"log"
)

func StartServer(addr string) {
	app := setupRouter()
   
	err := app.Listen(addr)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
