package app

const (
	ResponseStatusPassed = "passed"
	ResponseStatusFailed = "failed"
)

type ReportResponse struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Content string `json:"content"`
}
