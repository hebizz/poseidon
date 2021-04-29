package interfaces

type LogMsg struct {
	Uid         int    `json:"uid"`         //数据库自带uid
	Timestamp   int    `json:"timestamp"`   //时间戳
	EventType   string `json:"eventType"`   //类型(操作,告警)
	Description string `json:"description"` //描述
}

