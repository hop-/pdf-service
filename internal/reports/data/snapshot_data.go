package data

type SnapshotDataCompany struct {
	Name   string `json:"name"`
	Scores struct {
		Security   float32 `json:"security"`
		Compliance float32 `json:"compliance"`
		Financial  float32 `json:"financial"`
		Total      float32 `json:"total"`
	} `json:"scores"`
}

type SnapshotDataKpi struct {
	UnresolvedIssues                   int     `json:"unresolvedIssues"`
	ResolvedIssues                     int     `json:"resolvedIssues"`
	CriticalIssues                     int     `json:"criticalIssues"`
	VulnerabilityPatchMeantime         int     `json:"vulnerabilityPatchMeantime"`
	CriticalVulnerabilityPatchMeantime int     `json:"criticalVulnerabilityPatchMeantime"`
	BreachImpact                       int     `json:"breachImpact"`
	BreachLikelihood                   float32 `json:"breachLikelihood"`
}

type SnapshotDataIssues struct {
	CountsBySeverities SeverityCounts    `json:"countsBySeverities"`
	CountsByDimensions DimensionCounts   `json:"countsByDimensions"`
	CountByCategories  CategoriesCounts  `json:"countByCategories"`
	Compliances        CompliancesCounts `json:"compliances"`
	StatsBySeverity    []SeverityStat    `json:"statsBySeverity"`
}

type SnapshotData struct {
	Timestamp string              `json:"timestamp"`
	Company   SnapshotDataCompany `json:"company"`
	Kpi       SnapshotDataKpi     `json:"kpi"`
	Issues    SnapshotDataIssues  `json:"issues"`
}
