package crawl

import (
	"SHUVolunteer/model/activity"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func fetchActivityDetail(client http.Client, ctlId int, viewState string, viewStateGenerator string) (activity.Activity, error) {
	ctlIdString := fmt.Sprintf("ctl%02d", ctlId)
	response, _ := postFormWithUA(&client, "http://202.120.127.129/ShuLVMS/Volunteer/VActivityApply.aspx", url.Values{
		"ScriptManager1":       []string{"UpdatePanel1|GvActivityInfo$" + ctlIdString + "$LbtJionActivity"},
		"__EVENTTARGET":        []string{"GvActivityInfo$" + ctlIdString + "$LbtJionActivity"},
		"__EVENTARGUMENT":      []string{""},
		"__VIEWSTATE":          []string{viewState},
		"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
		"TxtActivityName":      []string{""},
		"DdlActivityType":      []string{"18"},
		"__ASYNCPOST":          []string{"true"},
	})
	responseBody, _ := ioutil.ReadAll(response.Body)
	idRegExp, _ := regexp.Compile("VATApply\\.aspx\\?Id=\\d+")
	regexMatchResult := idRegExp.Find(responseBody)
	if len(regexMatchResult) < 17 {
		return activity.Activity{}, errors.New("not valid activity")
	}
	id := string(regexMatchResult[17:])
	response, _ = client.Get("http://202.120.127.129/ShuLVMS/VandTandAInfo/ActivityInfo.aspx?Id=" + id)
	pageDoc, _ := goquery.NewDocumentFromReader(response.Body)
	title := pageDoc.Find("#LblActivityName").Text()
	leader := pageDoc.Find("#LblLeaderName").Text()
	address := pageDoc.Find("#LblActivityLocation").Text()
	startTime, _ := time.Parse("2006-1-02", pageDoc.Find("#LblActivityTime").Text())
	endTime, _ := time.Parse("2006-1-02", pageDoc.Find("#lblactivitytimeend").Text())
	return activity.Activity{
		Id:        id,
		Title:     title,
		Address:   address,
		Leader:    leader,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

func FetchAllActivities(client http.Client) []activity.Activity {
	activityList, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/VActivityApply.aspx")
	pageDoc, _ := goquery.NewDocumentFromReader(activityList.Body)
	viewState, _ := pageDoc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewStateGenerator, _ := pageDoc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
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
			name := tds[1]
			alreadyExists, err := activity.GetByName(name)
			if err == nil {
				result = append(result, alreadyExists)
				return ""
			}
			newActivity, err := fetchActivityDetail(client, i+1, viewState, viewStateGenerator)
			if err != nil {
				return ""
			}
			activity.Save(newActivity)
			result = append(result, newActivity)
			return ""
		}
	})
	return result
}

func FetchParticipatingActivities(client http.Client) []activity.Activity {
	allActivitiesResponse, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/ActivitySate.aspx")
	pageDoc, _ := goquery.NewDocumentFromReader(allActivitiesResponse.Body)
	names := pageDoc.Find("table table table table table tr").Map(func(i int, selection *goquery.Selection) string {
		return selection.Find("td:first-of-type").Text()
	})[1:]
	var activities []activity.Activity
	for _, name := range names {
		activityObject, err := activity.GetByName(name)
		if err != nil {
			FetchAllActivities(client)
			return FetchParticipatingActivities(client)
		}
		activities = append(activities, activityObject)
	}
	return activities
}

func ApplyActivity(client http.Client, activityId string) {
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

func ResignActivity(client http.Client, activityId string) {
	activityObject, _ := activity.Get(activityId)
	allActivitiesResponse, _ := client.Get("http://202.120.127.129/ShuLVMS/Volunteer/ActivitySate.aspx")
	pageDoc, _ := goquery.NewDocumentFromReader(allActivitiesResponse.Body)
	ctlString := strings.Join(pageDoc.Find("table table table table table tr").Map(func(i int, selection *goquery.Selection) string {
		if i == 0 {
			return ""
		} else {
			tds := selection.Find("td").Map(func(_ int, selection *goquery.Selection) string {
				return selection.Text()
			})
			if tds[0] == activityObject.Title {
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
