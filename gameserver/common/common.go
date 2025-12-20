package common

type TaskId interface {
	GenTaskId() string
}

type Request struct {
	UID       string                 `json:"uid"`
	RequestId int                    `json:"requestId"`
	Parma     map[string]interface{} `json:"parma"`
}

func (r Request) GenTaskId() string {
	if r.UID == "" {
		panic("uid is empty")
	}
	return r.UID
}
