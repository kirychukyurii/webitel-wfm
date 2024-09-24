package logger

import "encoding/json"

func (a Action) String() string {
	return string(a)
}

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type RequiredFields struct {
	UserId     int64  `json:"userId,omitempty"`
	UserIp     string `json:"userIp,omitempty"`
	Action     Action `json:"action,omitempty"`
	Date       int64  `json:"date,omitempty"`
	DomainId   int64
	ObjectName string
}

type Record struct {
	Id       int64  `json:"id,omitempty"`
	NewState []byte `json:"newState,omitempty"`
}

type Message struct {
	Records        []Record `json:"records,omitempty"`
	RequiredFields `json:"requiredFields"`
}

func (m *Message) ToJson() []byte {
	data, _ := json.Marshal(m)

	return data
}
