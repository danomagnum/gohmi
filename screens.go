package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

var screen_templates *template.Template

func parse_screens() {
	t, err := template.ParseGlob("./screens/*.html")
	if err != nil {
		slog.Error("problem parsing templates: %w", err)
		return
	}
	screen_templates = t
}

type ScreenData struct {
	Title string
}

func api_view(w http.ResponseWriter, req *http.Request) {
	screen_name := req.PathValue("screen")

	full_screen_name := fmt.Sprintf("%s.html", screen_name)

	if screen_name == "" {
		// TODO: go to a home page or something
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t := screen_templates.Lookup(full_screen_name)
	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dat := ScreenData{
		Title: screen_name,
	}

	w.Header().Set("Content-Type", "text/html")
	err := t.Execute(w, dat)
	if err != nil {
		slog.Error("problem with %s: %w", screen_name, err)
	}
}
