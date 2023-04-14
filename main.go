package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	AddRoutes(router)

	port := ResolvePort()
	log.Println("will server on port:", port)

	http.ListenAndServe(":"+port, router)
}
