package handler

import (
	"SHUVolunteer/model/student"
	"SHUVolunteer/service/crawl"
	"SHUVolunteer/service/token"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
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
	_, _ = w.Write(outputJSON)
	log.Printf("%s logged in", input.Username)
}

func getStudent(r *http.Request) (student.Student, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return student.Student{}, errors.New("no Authorization header given")
	}
	tokenString = tokenString[7:]
	return token.GetStudent(tokenString)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
