package main

import (
	"fmt"
	"groupie-tracker/server"
	"log"
	"net/http"
)

func main() {
	fmt.Println("http://localhost:8080/")
	fs := http.FileServer(http.Dir("./ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", server.MainPage)
	http.HandleFunc("/artists/", server.InfoAboutArtist)
	http.HandleFunc("/search/", server.SearchHandler)
	http.HandleFunc("/filter/", server.FilterhHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
