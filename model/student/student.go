package student

import (
	"SHUVolunteer/infrastructure"
	"encoding/json"
	"net/http"
)

type Student struct {
	Id      string
	Cookies []*http.Cookie
}

func marshal(student Student) []byte {
	marshaled, _ := json.Marshal(student)
	return marshaled
}

func unmarshal(binaryData []byte) Student {
	result := Student{}
	_ = json.Unmarshal(binaryData, &result)
	return result
}

func Save(student Student) {
	infrastructure.Redis.Set("Student_"+student.Id, marshal(student), 0)
}

func Get(id string) (Student, error) {
	binaryData, err := infrastructure.Redis.Get("Student_" + id).Result()
	return unmarshal([]byte(binaryData)), err
}
