package handler

import (
	"SHUVolunteer/service/crawl"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func ApplyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ActivityId string `json:"id"`
	}
	user, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	jar, _ := cookiejar.New(nil)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	jar.SetCookies(urlObject, user.Cookies)
	client := http.Client{Jar: jar}
	crawl.ApplyActivity(client, input.ActivityId)
}

func ResignHandler(w http.ResponseWriter, r *http.Request) {
	activityId := r.URL.Query().Get("activity_id")
	user, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	jar, _ := cookiejar.New(nil)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	jar.SetCookies(urlObject, user.Cookies)
	client := http.Client{Jar: jar}
	crawl.ResignActivity(client, activityId)
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
	// Go is kidding
	// @see https://github.com/golang/go/issues/26866
	if len(activities) == 0 {
		_, _ = w.Write([]byte("[]"))
	} else {
		response, _ := json.Marshal(activities)
		_, _ = w.Write(response)
	}
}

func ParticipatingActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	jar, _ := cookiejar.New(nil)
	urlObject, _ := url.Parse("http://202.120.127.129/")
	jar.SetCookies(urlObject, user.Cookies)
	client := http.Client{Jar: jar}
	activities := crawl.FetchParticipatingActivities(client)
	// Go is kidding
	// @see https://github.com/golang/go/issues/26866
	if len(activities) == 0 {
		_, _ = w.Write([]byte("[]"))
	} else {
		response, _ := json.Marshal(activities)
		_, _ = w.Write(response)
	}
}
func VolunteerActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("participating") == "true" {
			ParticipatingActivitiesHandler(w, r)
		} else {
			ActivityListHandler(w, r)
		}
	case "POST":
		ApplyHandler(w, r)
	case "DELETE":
		ResignHandler(w, r)
	}
}
