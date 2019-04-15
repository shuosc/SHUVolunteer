package crawl

import (
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

func Login(username string, password string) http.Client {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	pageResponse, _ := client.Get("http://202.120.127.129/Shulvms/Login.aspx")
	verificationCode := getVerificationCode(&client)
	x := strconv.Itoa(int(rand.Uint32()%20 + 1))
	y := strconv.Itoa(int(rand.Uint32()%20 + 1))
	pageDoc, _ := goquery.NewDocumentFromReader(pageResponse.Body)
	viewState, _ := pageDoc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewStateGenerator, _ := pageDoc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	_, _ = postFormWithUA(&client, "http://202.120.127.129/Shulvms/Login.aspx", url.Values{
		"ScriptManager1":       []string{"UpdatePanel1|bt_login"},
		"__EVENTTARGET":        []string{""},
		"__EVENTARGUMENT":      []string{""},
		"tb_user":              []string{username},
		"tb_pass":              []string{password},
		"yzm":                  []string{verificationCode},
		"__VIEWSTATE":          []string{viewState},
		"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
		"__ASYNCPOST":          []string{"true"},
		"bt_login.x":           []string{x},
		"bt_login.y":           []string{y},
	})
	return client
}
