// File: main.go
// Creation: Wed May 29 16:35:03 2024
// Time-stamp: <2024-06-07 14:53:52>
// Copyright (C): 2024 Pierre Lecocq

package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	port := 8080
	db := NewDatabase("dailtasks.db")

	log.Printf("Starting server on port %d\n", port)
	log.Printf("Please visit http://localhost:%d/\n", port)

	StartServer(port, db)
}
