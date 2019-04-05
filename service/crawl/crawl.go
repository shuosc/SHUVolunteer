package crawl

import (
	"SHUVolunteer/model/activity"
	"SHUVolunteer/service"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getVerificationCode(client *http.Client) string {
	verificationCodeResponse, _ := client.Get("http://202.120.127.129/Shulvms/GetPic.aspx")
	verificationCodeBody, _ := ioutil.ReadAll(verificationCodeResponse.Body)
	verificationCode, _ := service.Ocr(verificationCodeBody)
	return verificationCode
}

func postFormWithUA(client *http.Client, url string, form url.Values) (*http.Response, error) {
	const UA string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"
	req, _ := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", UA)
	return client.Do(req)
}

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

func FetchAllActivities(client http.Client) []activity.Activity {
	activityList, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/VActivityApply.aspx")
	pageDoc, _ := goquery.NewDocumentFromReader(activityList.Body)
	var result []activity.Activity
	pageDoc.Find("#UpdatePanel1 table table table table table tr").Map(func(i int, selection *goquery.Selection) string {
		if i == 0 {
			return ""
		} else {
			tds := selection.Find("td").Map(func(_ int, selection *goquery.Selection) string {
				return selection.Text()
			})
			if len(tds) != 7 {
				return ""
			}
			date, _ := time.Parse("2006-01-02", tds[0])
			result = append(result, activity.Activity{
				Date: date,
				Name: tds[1],
				Team: tds[2],
			})
			return ""
		}
	})
	return result
}

func fetchActivityID(client http.Client, allActivities []activity.Activity, activityName string) string {
	activityList, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/VActivityApply.aspx")
	pageDoc, _ := goquery.NewDocumentFromReader(activityList.Body)
	viewState, _ := pageDoc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewStateGenerator, _ := pageDoc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	for index, activity := range allActivities {
		if activity.Name == activityName {
			ctlId := fmt.Sprintf("ctl%02d", index+2)
			response, _ := postFormWithUA(&client, "http://202.120.127.129/ShuLVMS/Volunteer/VActivityApply.aspx", url.Values{
				"ScriptManager1":       []string{"UpdatePanel1|GvActivityInfo$" + ctlId + "$LbtJionActivity"},
				"__EVENTTARGET":        []string{"GvActivityInfo$" + ctlId + "$LbtJionActivity"},
				"__EVENTARGUMENT":      []string{""},
				"__VIEWSTATE":          []string{viewState},
				"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
				"TxtActivityName":      []string{""},
				"DdlActivityType":      []string{"18"},
				"__ASYNCPOST":          []string{"true"},
			})
			responseBody, _ := ioutil.ReadAll(response.Body)
			idRegExp, _ := regexp.Compile("VATApply\\.aspx\\?Id=\\d+")
			return string(idRegExp.Find(responseBody)[17:])
		}
	}
	return ""
}

func applyActivityWithId(client http.Client, activityId string) {
	pageResponse, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/VATApply.aspx?Id=" + activityId)
	pageDoc, _ := goquery.NewDocumentFromReader(pageResponse.Body)
	viewState, _ := pageDoc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewStateGenerator, _ := pageDoc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	activityTime := pageDoc.Find("#lblactivitytime").Text()
	_, _ = postFormWithUA(&client, "http://202.120.127.129/ShuLVMS/Volunteer/VATApply.aspx?Id="+activityId, url.Values{
		"ScriptManager1":       []string{"UpdatePanel1|BtnSubmit"},
		"__EVENTTARGET":        []string{""},
		"__EVENTARGUMENT":      []string{""},
		"__VIEWSTATE":          []string{viewState},
		"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
		"DropDownList1":        []string{"1"},
		"txtbeizhu":            []string{""},
		"txtActivityTime":      []string{activityTime},
		"txtActivityTimeEnd":   []string{activityTime},
		"__ASYNCPOST":          []string{"true"},
		"BtnSubmit":            []string{"我要加入"},
	})
}

func ApplyActivity(client http.Client, activityName string) {
	activities := FetchAllActivities(client)
	applyActivityWithId(client, fetchActivityID(client, activities, activityName))
}

func ResignActivity(client http.Client, activityName string) {
	allActivitiesResponse, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/ActivitySate.aspx")
	pageDoc, _ := goquery.NewDocumentFromReader(allActivitiesResponse.Body)
	ctlString := strings.Join(pageDoc.Find("table table table table table tr").Map(func(i int, selection *goquery.Selection) string {
		if i == 0 {
			return ""
		} else {
			tds := selection.Find("td").Map(func(_ int, selection *goquery.Selection) string {
				return selection.Text()
			})
			if tds[0] == activityName {
				return fmt.Sprintf("ctl%02d", i+1)
			} else {
				return ""
			}
		}
	}), "")
	viewState, _ := pageDoc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewStateGenerator, _ := pageDoc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	_, _ = postFormWithUA(&client, "http://202.120.127.129/ShuLVMS/Volunteer/ActivitySate.aspx", url.Values{
		"ScriptManager1":       []string{"UpdatePanel1|GvTeamAvailable$" + ctlString + "$LbtView"},
		"__EVENTTARGET":        []string{"GvTeamAvailable$" + ctlString + "$LbtView"},
		"__EVENTARGUMENT":      []string{""},
		"__VIEWSTATE":          []string{viewState},
		"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
		"TxtSearch":            []string{""},
		"__ASYNCPOST":          []string{"true"},
	})
}
