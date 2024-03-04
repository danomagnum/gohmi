package logix

import (
	"html/template"
	"net/http"
)

func (drv *LogixDriver) Display() template.HTML {
	return ""
}

func (drv *LogixDriver) Edit(req *http.Request) template.HTML {
	return ""
}
