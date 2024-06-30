package reports

type Engine interface {
	Generate(pages [][]byte, styles, headerHtml, footerHtml *string) ([]byte, error)
}
