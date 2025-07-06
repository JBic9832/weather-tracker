package main

import "log"

func main() {
	server := NewServer(":8000")
	if err := server.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
