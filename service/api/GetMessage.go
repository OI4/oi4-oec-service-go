package api

type GetMessage struct {
	MessageId      string  `json:"MessageId"`
	MessageType    string  `json:"MessageType"`
	PublisherId    string  `json:"PublisherId"`
	DataSetClassId string  `json:"DataSetClassId"`
	CorrelationId  *string `json:"CorrelationId"`
	Messages       []any   `json:"Messages"`
}
