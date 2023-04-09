package main

import (
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/server_interface/rest_server"
)

func main() {
	rest_server.StartServer(":8080")
}
