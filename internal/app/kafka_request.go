package app

import "encoding/json"

type KafkaRequest struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Data any    `json:"data"`
}

func (r *KafkaRequest) JsonData() ([]byte, error) {
	return json.Marshal(r.Data)
}
