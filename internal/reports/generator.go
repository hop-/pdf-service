package reports

import (
	"encoding/base64"
)

type ReportGenerator struct {
	template     *Template
	pdfGenerator Engine
}

func NewReportGenerator(templateName string, engineName string) (*ReportGenerator, error) {
	t, err := NewTemplate(templateName)
	if err != nil {
		return nil, err
	}

	var pdfg Engine

	switch engineName {
	case "wkhtmltopdf":
		pdfg, err = NewWkhtpEngine()
		if err != nil {
			return nil, err
		}
	default:
		pdfg = NewChrdpEngine()
	}

	return &ReportGenerator{t, pdfg}, nil
}

func (r *ReportGenerator) GenerateBase64(data []byte) (string, error) {
	pdf, err := r.Generate(data)

	return base64.StdEncoding.EncodeToString(pdf), err
}

func (r *ReportGenerator) HasCoverPage() bool {
	return r.template.coverTmpl != nil
}

func (r *ReportGenerator) Generate(data []byte) ([]byte, error) {
	content := [][]byte{}
	if r.HasCoverPage() {
		coverHtmlContent, err := r.template.GenerateCoverPage(data)
		if err != nil {
			return []byte{}, err
		}

		content = append(content, coverHtmlContent)
	}
	// TODO: add table of content

	htmlContent, err := r.template.GenerateContent(data)
	if err != nil {
		return []byte{}, err
	}

	content = append(content, htmlContent)
	styles := r.template.GetReportStyles()
	headerHtml := r.template.GetHeaderHtml()
	footerHtml := r.template.GetFooterHtml()

	return r.pdfGenerator.Generate(content, styles, headerHtml, footerHtml)
}
