package reports

import (
	"bytes"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type WkhtpEngine struct {
	pdfGenerator *wkhtmltopdf.PDFGenerator
}

func NewWkhtpEngine() (*WkhtpEngine, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	//pdfg.MarginTop.Set(0)
	//pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	return &WkhtpEngine{pdfg}, nil
}

func (e *WkhtpEngine) Generate(pages [][]byte, styles, headerHtml, footerHtml *string) ([]byte, error) {
	e.setReportAllPages(pages, styles, headerHtml, footerHtml, true)

	err := e.pdfGenerator.Create()
	if err != nil {
		return nil, err
	}

	return e.pdfGenerator.Bytes(), nil
}

func (r *WkhtpEngine) setReportAllPages(pages [][]byte, styles, headerHtml, footerHtml *string, pageCount bool) {
	buf := new(bytes.Buffer)
	for i := range pages {
		buf.Write(pages[i])
		buf.WriteString(`<p style="page-break-before: always"></p>`)
	}

	p := wkhtmltopdf.NewPageReader(buf)

	// Make sure any injections are not allowed
	p.EnableLocalFileAccess.Set(true)
	p.DisableExternalLinks.Set(false)
	p.DisableSmartShrinking.Set(true)

	if styles != nil {
		p.UserStyleSheet.Set(*styles)
	}

	if headerHtml != nil {
		p.HeaderHTML.Set(*headerHtml)
	}

	if footerHtml != nil {
		p.FooterHTML.Set(*footerHtml)
	}

	if pageCount {
		p.FooterRight.Set("[page]/[topage]    ")
	}

	r.pdfGenerator.AddPage(p)
}
