package model

type RabbitMQMsg struct {
	MsgId     string  `json:"msgId,omitempty"`
	CommandId string  `json:"commandId,omitempty"`
	Version   string  `json:"version,omitempty"`
	Topic     string  `json:"topic,omitempty"`
	Code      int     `json:"code,omitempty"`
	Message   string  `json:"message,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
	Params    []Param `json:"params,omitempty"`
}

type Param struct {
	Id        string `json:"id,omitempty"`
	Value     string `json:"value,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}
