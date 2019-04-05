package main

import (
	"SHUVolunteer/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/apply", handler.ApplyHandler)
	http.HandleFunc("/resign", handler.ResignHandler)
	http.HandleFunc("/activities", handler.ActivityListHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
