package handler

import (
	"SHUVolunteer/model/student"
	"SHUVolunteer/service/crawl"
	"SHUVolunteer/service/token"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func getStudent(r *http.Request) (student.Student, error) {
	tokenString := r.Header.Get("Authorization")[7:]
	return token.GetStudent(tokenString)
}

func ApplyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ActivityName string `json:"activity_name"`
	}
	user, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	jar, _ := cookiejar.New(nil)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	jar.SetCookies(urlObject, user.Cookies)
	client := http.Client{Jar: jar}
	crawl.ApplyActivity(client, input.ActivityName)
}

func ResignHandler(w http.ResponseWriter, r *http.Request) {
	activityName := r.URL.Query().Get("activity_name")
	user, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	jar, _ := cookiejar.New(nil)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	jar.SetCookies(urlObject, user.Cookies)
	client := http.Client{Jar: jar}
	crawl.ResignActivity(client, activityName)
}

func ActivityListHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	jar, _ := cookiejar.New(nil)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	jar.SetCookies(urlObject, user.Cookies)
	client := http.Client{Jar: jar}
	activities := crawl.FetchAllActivities(client)
	response, _ := json.Marshal(activities)
	w.Write(response)
}
