package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	stmp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	// format := time.Time(t).Format("2022-10-10 12:45:12")
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
