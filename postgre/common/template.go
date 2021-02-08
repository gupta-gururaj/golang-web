package common

import (
	"go/build"
	"text/template"
)

//GetTemplate is ...
func GetTemplate() template.Template {
	var route = build.Default.GOPATH + "/src/postgre/templates/*"
	var tmpl = template.Must(template.ParseGlob(route))
	return *tmpl
}
