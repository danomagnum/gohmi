package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func api_write(w http.ResponseWriter, req *http.Request) {
	driver_name := req.PathValue("driver")
	tag_name := req.PathValue("tag")

	defer req.Body.Close()

	if driver_name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("No driver name specified."))
		if err != nil {
			slog.Error("failed to write bad driver msg to %s: %w", req.RemoteAddr, err)
		}
		return
	}

	if tag_name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("No tag name specified."))
		if err != nil {
			slog.Error("failed to write bad tagname msg to %s: %w", req.RemoteAddr, err)
		}
		return
	}

	driver, ok := drivers[driver_name]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(fmt.Sprintf("driver %s was not found", driver_name)))
		if err != nil {
			slog.Error("failed to write unknown driver name msg to %s: %w", req.RemoteAddr, err)
		}
		return
	}

	jdec := json.NewDecoder(req.Body)
	var val any

	jdec.Decode(&val)

	err := driver.Write(tag_name, val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("error writing %s/%s = %v: %v", driver_name, tag_name, val, err)))
		if err != nil {
			slog.Error("failed to write unknown driver name msg to %s: %w", req.RemoteAddr, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

}
