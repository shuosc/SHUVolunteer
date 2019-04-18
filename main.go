package main

import (
	"SHUVolunteer/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/ping", handler.TestHandler)
	http.HandleFunc("/activities", handler.VolunteerActivitiesHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
