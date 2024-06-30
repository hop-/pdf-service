package data

type SeverityCounts struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type DimensionCounts struct {
	Api               int `json:"API"`
	Email             int `json:"Email"`
	Network           int `json:"Network"`
	Patching          int `json:"Patching"`
	WebApplication    int `json:"Web Application"`
	Reputation        int `json:"Reputation"`
	Dns               int `json:"DNS"`
	MalwareRansomware int `json:"Malwar and Ransomware"`
	PhishingDataLeak  int `json:"Phishing and Data Leak"`
}

type CategoriesCounts struct {
	Api               SeverityCounts `json:"API"`
	Email             SeverityCounts `json:"Email"`
	Network           SeverityCounts `json:"Network"`
	Patching          SeverityCounts `json:"Patching"`
	WebApplication    SeverityCounts `json:"Web Application"`
	Reputation        SeverityCounts `json:"Reputation"`
	Dns               SeverityCounts `json:"DNS"`
	MalwareRansomware SeverityCounts `json:"Malwar and Ransomware"`
	PhishingDataLeak  SeverityCounts `json:"Phishing and Data Leak"`
}

type CompliancesCounts struct {
	Nist80053 int `json:"Nist 800-53"`
	Iso27001  int `json:"ISO 27001"`
	// TODO
}

type SeverityStat struct {
	Date   string         `json:"date"`
	Counts SeverityCounts `json:"counts"`
}
