package handler

import (
	"SHUVolunteer/model/student"
	"SHUVolunteer/service/crawl"
	"SHUVolunteer/service/token"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	loginStudentClient := crawl.Login(input.Username, input.Password)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	cookies := loginStudentClient.Jar.Cookies(urlObject)
	studentObject := student.Student{
		Id:      input.Username,
		Cookies: cookies,
	}
	var output struct {
		Token string `json:"token"`
	}
	output.Token = token.GenerateJWT(studentObject)
	student.Save(studentObject)
	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}

func VolunteerActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ActivityListHandler(w, r)
	case "POST":
		ApplyHandler(w, r)
	case "DELETE":
		ResignHandler(w, r)
	}
}
