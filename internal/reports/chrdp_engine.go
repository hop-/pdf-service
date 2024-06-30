package reports

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ChrdpEngine struct{}

func NewChrdpEngine() *ChrdpEngine {
	return &ChrdpEngine{}
}

func (e *ChrdpEngine) Generate(pages [][]byte, styles, headerHtml, footerHtml *string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	file, err := os.CreateTemp("./", "*.html")
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	for i := range pages {
		file.Write(pages[i])
		// TODO make it optional from templaterc.json configs
		// file.WriteString(`<p style="page-break-before: always"></p>`)
	}

	filePath, err := filepath.Abs(file.Name())
	if err != nil {
		return nil, err
	}

	var buf []byte

	// TODO: Add header footer styles
	err = chromedp.Run(ctx,
		chromedp.Navigate((&url.URL{Scheme: "file", Path: filePath}).String()),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithDisplayHeaderFooter(true).
				Do(ctx)
			if err != nil {
				return err
			}

			return nil
		}),
	)

	return buf, err
}
