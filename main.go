package main

import (
	"SHUVolunteer/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/volunteer-activities", handler.VolunteerActivitiesHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
