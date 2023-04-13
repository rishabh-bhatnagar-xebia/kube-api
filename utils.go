package main

import (
	"log"
	"os"
)

func resolvePort() string {
	port := PORT
	if len(os.Args) > 1 {
		port = os.Args[1]
		log.Println("got an inline value for the port:", port)
	}
	return port
}
