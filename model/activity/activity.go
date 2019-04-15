package activity

import (
	"SHUVolunteer/infrastructure"
	"encoding/json"
	"time"
)

type Activity struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Leader    string    `json:"leader"`
	Address   string    `json:"address"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func marshal(activity Activity) []byte {
	marshaled, _ := json.Marshal(activity)
	return marshaled
}

func unmarshal(binaryData []byte) Activity {
	result := Activity{}
	_ = json.Unmarshal(binaryData, &result)
	return result
}

func Save(activity Activity) {
	infrastructure.Redis.Set("Activity_"+activity.Id+"_"+activity.Title, marshal(activity), 0)
}

func Get(id string) (Activity, error) {
	key, _ := infrastructure.Redis.Keys("Activity_" + id + "_*").Result()
	binaryData, err := infrastructure.Redis.Get(key[0]).Result()
	return unmarshal([]byte(binaryData)), err
}

func GetByName(name string) (Activity, error) {
	key, _ := infrastructure.Redis.Keys("Activity_*_" + name).Result()
	binaryData, err := infrastructure.Redis.Get(key[0]).Result()
	return unmarshal([]byte(binaryData)), err
}
