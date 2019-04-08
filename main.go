package main

import (
	"SHUVolunteer/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/volunteer-activities", handler.VolunteerActivitiesHandle)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
