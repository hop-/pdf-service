package data

type DiffValue[T any] struct {
	Current  T `json:"current"`
	Previous T `json:"previous"`
	Diff     struct {
		Amount    T       `json:"amount"`
		Percent   float32 `json:"percent"`
		Direction string  `json:"direction"`
		Status    string  `json:"status"`
	} `json:"diff"`
}

type DiffDataCompany struct {
	Name   string `json:"name"`
	Scores struct {
		Security   DiffValue[float32] `json:"security"`
		Compliance DiffValue[float32] `json:"compliance"`
		Financial  DiffValue[float32] `json:"financial"`
		Total      DiffValue[float32] `json:"total"`
	} `json:"scores"`
}

type DiffDataAssets struct {
	DomainsAndSubdomains DiffValue[int] `json:"domainsAndSubdomains"`
	Websites             DiffValue[int] `json:"websites"`
	CloudServices        DiffValue[int] `json:"cloudServices"`
	Ips                  DiffValue[int] `json:"ips"`
	Ports                DiffValue[int] `json:"ports"`
}

type DiffDataKpi struct {
	UnresolvedIssues                   DiffValue[int]     `json:"unresolvedIssues"`
	ResolvedIssues                     DiffValue[int]     `json:"resolvedIssues"`
	CriticalIssues                     DiffValue[int]     `json:"criticalIssues"`
	VulnerabilityPatchMeantime         DiffValue[int]     `json:"vulnerabilityPatchMeantime"`
	CriticalVulnerabilityPatchMeantime DiffValue[int]     `json:"criticalVulnerabilityPatchMeantime"`
	BreachImpact                       DiffValue[int]     `json:"breachImpact"`
	BreachLikelihood                   DiffValue[float32] `json:"breachLikelihood"`
}

type DiffDataIssues struct {
	BySeverities struct {
		Low      DiffValue[int] `json:"low"`
		Medium   DiffValue[int] `json:"medium"`
		High     DiffValue[int] `json:"high"`
		Critical DiffValue[int] `json:"critical"`
	} `json:"bySeverities"`
	ByDimensions struct {
		Api               DiffValue[int] `json:"API"`
		Email             DiffValue[int] `json:"Email"`
		Network           DiffValue[int] `json:"Network"`
		Patching          DiffValue[int] `json:"Patching"`
		WebApplication    DiffValue[int] `json:"Web Application"`
		Reputation        DiffValue[int] `json:"Reputation"`
		Dns               DiffValue[int] `json:"DNS"`
		MalwareRansomware DiffValue[int] `json:"Malwar and Ransomware"`
		PhishingDataLeak  DiffValue[int] `json:"Phishing and Data Leak"`
	} `json:"byDimensions"`
	StatsBySeverity      []SeverityStat `json:"statsBySeverity"`
	ComplianceViolations map[string]DiffValue[int] `json:"complianceViolations"`
	New []struct {
		HssId          string   `json:"hssId"`
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		Severity       string   `json:"severity"`
		Dimension      string   `json:"dimension"`
		BusinessRisk   string   `json:"businessRisk"`
		ImpactedAssets []string `json:"impactedAssets"`
	} `json:"new"`
}

type DiffData struct {
	Timestamp         string          `json:"timestamp"`
	PreviousTimestamp string          `json:"previousTimestamp"`
	Company           DiffDataCompany `json:"company"`
	FinancialRisk     DiffValue[int]  `json:"financialRisk"`
	Kpi               DiffDataKpi     `json:"kpi"`
	Assets            DiffDataAssets  `json:"assets"`
	Issues            DiffDataIssues  `json:"issues"`
}
