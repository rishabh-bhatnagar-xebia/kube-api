package main

import (
	"encoding/json"
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

func wrap(err error) string {
	content, _ := json.Marshal(struct {
		Error string `json: err`
	}{
		Error: err.Error(),
	})
	return string(content)
}
