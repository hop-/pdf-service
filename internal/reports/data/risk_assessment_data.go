package data

type RiskAssessmentDataCompany struct {
	Name      string `json:"name"`
	Industry  string `json:"industry"`
	Employees string `json:"employees"`
	Scores    struct {
		Security   float32 `json:"security"`
		Compliance float32 `json:"compliance"`
		Financial  float32 `json:"financial"`
		Total      float32 `json:"total"`
	} `json:"scores"`
}

type RiskAssessmentDataAssets struct {
	MainDomains    []string `json:"mainDomains"`
	RelatedDomains []string `json:"relatedDomains"`
	Counts         struct {
		DomainsAndSubdomains    int `json:"domainsAndSubdomains"`
		WebSitesAndCertificates int `json:"websitesAndCertificates"`
		IpsAndPorts             int `json:"ipsAndPorts"`
		CloudsAndNetworks       int `json:"cloudsAndNetworks"`
		EmailsAndPhishings      int `json:"emailsAndPhishings"`
	} `json:"counts"`
	DomainsAndSubDomains []string `json:"domainsAndSubdomains"`
	Websites             []string `json:"websites"`
	Ips                  []string `json:"ips"`
	IpBlocks             []struct {
		Name    string `json:"name"`
		Owner   string `json:"owner"`
		Asn     string `json:"asn"`
		AsnType string `json:"asnType"`
	} `json:"ipBlocks"`
	Ports []struct {
		Port        int    `json:"port"`
		DomainCount int    `json:"domainCount"`
		IpCount     int    `json:"ipCount"`
		Severity    string `json:"severity"`
	} `json:"ports"`
}

type RiskAssessmentDataImpersonation struct {
	Domain     string `json:"domain"`
	HttpBanner string `json:"httpBanner"`
	Dns        struct {
		A    []string `json:"a"`
		Aaaa []string `json:"aaaa"`
		Mx   []string `json:"mx"`
		Ns   []string `json:"ns"`
	} `json:"dns"`
	Whois struct {
		CreatedAt string `json:"createdAt"`
		Registrar string `json:"registrar"`
	} `json:"whois"`
	FirstSeen   string `json:"firstSeen"`
	LastUpdated string `json:"lastUpdated"`
}

type RiskAssessmentDataVulnerableAsset struct {
	Domain     string         `json:"domain"`
	Severities SeverityCounts `json:"severities"`
}

type RiskAssessmentDataWeakness struct {
	Name          string   `json:"name"`
	Severity      string   `json:"severity"`
	BusinessRisk  string   `json:"businessRisk"`
	BaseScore     float32  `json:"baseScore"`
	AttackVectors []string `json:"attackVectors"`
}

type RiskAssessmentDataComplianceViolationsCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type RiskAssessmentDataSecurityIssue struct {
	Name           string   `json:"name"`
	Seveirty       string   `json:"severity"`
	BusinessRisk   string   `json:"businessRisk"`
	Dimension      string   `json:"dimension"`
	Description    string   `json:"description"`
	BaseScore      float32  `json:"baseScore"`
	Recommendation string   `json:"recommendation"`
	References     []string `json:"references"`
	AttackVectors  []string `json:"attackVectors"`
	ImpactedAssets []struct {
		Name      string `json:"name"`
		FirstSeen string `json:"firstSeen"`
		// TODO: add Proofs `json:"proofs"`
	} `json:"impactedAssets"`
}

type RiskAssessmentDataSecurityIssuesByDimension struct {
	Name               string                            `json:"name"`
	TotalCount         int                               `json:"totalCount"`
	CountsBySeverities SeverityCounts                    `json:"countsBySeverities"`
	Issues             []RiskAssessmentDataSecurityIssue `json:"issues"`
}

type RiskAssessmentData struct {
	Timestamp                  string                                        `json:"timestamp"`
	Company                    RiskAssessmentDataCompany                     `json:"company"`
	Assets                     RiskAssessmentDataAssets                      `json:"assets"`
	Impersonations             []RiskAssessmentDataImpersonation             `json:"impersonations"`
	TopVulnerableAssets        []RiskAssessmentDataVulnerableAsset           `json:"topVulnerableAssets"`
	TopAttacVectors            []string                                      `json:"topAttackVectors"`
	TopWeaknesses              []RiskAssessmentDataWeakness                  `json:"topWeaknesses"`
	SecurityIssuesCounts       SeverityCounts                                `json:"securityIssuesCounts"`
	ComplianceViolationsCounts []RiskAssessmentDataComplianceViolationsCount `json:"complianceViolationsCounts"`
	SecurityIssuesByDimensions []RiskAssessmentDataSecurityIssuesByDimension `json:"securityIssuesByDimensions"`
}
