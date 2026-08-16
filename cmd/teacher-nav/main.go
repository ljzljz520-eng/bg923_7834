package main

import (
	"flag"
	"log"
	"net/http"

	"example.com/teacher-resource-navigation/internal/navigation"
)

func main() {
	address := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	service := navigation.NewServiceWithFixtures()
	server := &http.Server{Addr: *address, Handler: navigation.NewHandler(service)}
	log.Printf("teacher resource navigation listening on %s", *address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
