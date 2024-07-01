package app

const (
	ResponseStatusPassed = "passed"
	ResponseStatusFailed = "failed"
)

type KafkaResponse struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Content string `json:"content"`
}
