package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func api_view(w http.ResponseWriter, req *http.Request) {
	screen_name := req.PathValue("screen")

	full_screen_name := filepath.Join(".", "screens", fmt.Sprintf("%s.html", screen_name))
	//full_screen_name := fmt.Sprintf(".screens/%s.html", screen_name)

	if screen_name == "" {
		// TODO: go to a home page or something
		w.WriteHeader(http.StatusNotFound)
		return
	}
	subtemplate, err := gonja.FromFile(full_screen_name)
	if err != nil || subtemplate == nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "could not get screen html file\n")
		io.WriteString(w, err.Error())
		return
	}

	page_data, err := subtemplate.ExecuteToString(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "could not get execute page\n")
		io.WriteString(w, err.Error())
		return
	}
	template, err := gonja.FromFile("./screens/main.html")
	if err != nil || subtemplate == nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "could not get main screen html file\n")
		io.WriteString(w, err.Error())
		return
	}

	dat := exec.NewContext(map[string]any{
		"Title": screen_name,
		"Body":  page_data,
	})
	final_page_data := strings.Builder{}
	err = template.Execute(&final_page_data, dat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "could not get execute main page with content\n")
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, final_page_data.String())
}
