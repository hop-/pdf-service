package app

import "encoding/json"

type ReportRequest struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Data any    `json:"data"`
}

func (r *ReportRequest) JsonData() ([]byte, error) {
	return json.Marshal(r.Data)
}
