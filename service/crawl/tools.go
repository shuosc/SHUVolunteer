package crawl

import (
	"SHUVolunteer/service/ocr"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
