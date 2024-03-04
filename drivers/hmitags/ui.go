package hmitags

import (
	"html/template"
	"net/http"
)

func (drv *TagStore) Display() template.HTML {
	return ""
}

func (drv *TagStore) Edit(req *http.Request) template.HTML {
	return ""
}
