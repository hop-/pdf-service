package reports

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/reports/data"
)

var reportTemplates map[string]*Template

type Data struct {
	Context *Context
	Data    any
}

type Config struct {
	Name      string   `json:"name"`
	Cover     string   `json:"cover"`
	Templates []string `json:"templates"`
	Style     string   `json:"style"`
	Header    string   `json:"header"`
	Footer    string   `json:"footer"`
}

type Template struct {
	name      string
	context   Context
	coverTmpl *template.Template
	tmpl      *template.Template
	header    string
	footer    string
	styles    string
}

const templatesPath string = "templates" // TODO: make configurable

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !(err != nil && errors.Is(err, os.ErrNotExist))
}

func readConfig(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)

	return &config, err
}

func initTemplates() error {
	reportTemplates = map[string]*Template{}

	templatesDirPath, err := filepath.Abs(templatesPath)
	if err != nil {
		return err
	}

	golog.Info("Scanning", templatesDirPath, "for templates")

	templatesDirFile, err := os.Open(templatesDirPath)
	if err != nil {
		return err
	}

	templatesDirContent, err := templatesDirFile.ReadDir(0)
	if err != nil {
		return err
	}

	for _, d := range templatesDirContent {
		if !d.IsDir() {
			continue
		}
		templateDirPath := filepath.Join(templatesDirPath, d.Name())
		templateRcFilePath := filepath.Join(templateDirPath, "templaterc.json")
		if !fileExists(templateRcFilePath) {
			golog.Info(d.Name(), "is not a template, skipping")
			continue
		}
		golog.Info("Found new template", d.Name())
		config, err := readConfig(templateRcFilePath)
		if err != nil {
			golog.Error("Failed to read templaterc file:", err.Error())
			continue
		}

		if config.Name == "" {
			golog.Error("Template name is unspecified, skipping")
			continue
		}

		// TODO: Check if report with template name is already exist

		// TODO: currently using general data
		reportTemplate, err := initTemplate(templateDirPath, config)
		if err != nil {
			return err
		}

		reportTemplates[config.Name] = reportTemplate
	}

	return nil
}

func initTemplate(templateDir string, config *Config) (*Template, error) {
	htmlTemplateFiles := []string{}
	for _, fn := range config.Templates {
		filePath := filepath.Join(templateDir, fn)
		if !fileExists(filePath) {
			golog.Warning("Unable to find template file", fn, "specified in templaterc")
			continue
		}
		htmlTemplateFiles = append(htmlTemplateFiles, filePath)
	}

	footerFilePath := ""
	if config.Footer != "" {
		filePath := filepath.Join(templateDir, config.Footer)
		if fileExists(filePath) {
			footerFilePath = filePath
		} else {
			golog.Warning("Unable to find footer file", config.Footer, "specified in templaterc")
		}
	}

	headerFilePath := ""
	if config.Header != "" {
		filePath := filepath.Join(templateDir, config.Header)
		if fileExists(filePath) {
			headerFilePath = filePath
		} else {
			golog.Warning("Unable to find header file", config.Header, "specified in templaterc")
		}
	}

	stylesFilePath := ""
	if config.Style != "" {
		filePath := filepath.Join(templateDir, config.Style)
		if fileExists(filePath) {
			stylesFilePath = filePath
		} else {
			golog.Warning("Unable to find styles file", config.Style, "specified in templaterc")
		}
	}

	tmpl, err := template.ParseFiles(htmlTemplateFiles...)
	if err != nil {
		return nil, err
	}

	var coverTmpl *template.Template = nil
	if config.Cover != "" {
		filePath := filepath.Join(templateDir, config.Cover)
		if fileExists(filePath) {
			coverTmpl, err = template.ParseFiles(filePath)
			if err != nil {
				return nil, err
			}
		} else {
			golog.Warning("Unable to find styles file", config.Style, "specified in templaterc")
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	t := Template{
		name:      config.Name,
		context:   NewContext(pwd, templateDir),
		coverTmpl: coverTmpl,
		tmpl:      tmpl,
		header:    headerFilePath,
		footer:    footerFilePath,
		styles:    stylesFilePath,
	}

	return &t, nil
}

func NewTemplate(name string) (*Template, error) {
	if len(reportTemplates) == 0 {
		err := initTemplates()
		if err != nil {
			return nil, err
		}
	}

	tmpl, ok := reportTemplates[name]
	if !ok {
		// TODO: send error with message
		return nil, fmt.Errorf("unknown template name %s", name)
	}

	return tmpl, nil
}

func (t *Template) GenerateCoverPage(marshaledData []byte) ([]byte, error) {
	var data data.GeneralData
	err := json.Unmarshal(marshaledData, &data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = t.coverTmpl.Execute(&buf, Data{Context: &t.context, Data: data})

	return buf.Bytes(), err
}

func (t *Template) GenerateContent(marshaledData []byte) ([]byte, error) {
	var data data.GeneralData
	err := json.Unmarshal(marshaledData, &data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = t.tmpl.Execute(&buf, Data{Context: &t.context, Data: data})

	return buf.Bytes(), err
}

func (t *Template) GetReportStyles() *string {
	if t.styles == "" {
		return nil
	}

	return &t.styles
}

func (t *Template) GetHeaderHtml() *string {
	if t.header == "" {
		return nil
	}

	return &t.header
}

func (t *Template) GetFooterHtml() *string {
	if t.footer == "" {
		return nil
	}

	return &t.footer
}
