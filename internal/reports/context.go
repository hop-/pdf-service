package reports

import (
	"html/template"
	"net/url"
)

type Context struct {
	Pwd     string
	PwdUri  template.URL
	Root    string
	RootUri template.URL
}

func NewContext(pwd string, root string) Context {
	return Context{
		Pwd:     pwd,
		PwdUri:  template.URL((&url.URL{Scheme: "file", Path: pwd}).String()),
		Root:    root,
		RootUri: template.URL((&url.URL{Scheme: "file", Path: root}).String()),
	}
}
