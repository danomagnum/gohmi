package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func api_read(w http.ResponseWriter, req *http.Request) {
	driver_name := req.PathValue("driver")
	tag_name := req.PathValue("tag")

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

	val, err := driver.Read(tag_name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("error reading %s/%s: %v", driver_name, tag_name, err)))
		if err != nil {
			slog.Error("failed to write unknown driver name msg to %s: %w", req.RemoteAddr, err)
		}
		return
	}

	jenc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	err = jenc.Encode(val)
	if err != nil {
		slog.Error("failed to write json encoded value for %s/%s to %s: %w", driver_name, tag_name, req.RemoteAddr, err)
	}
}

func api_read_multi(w http.ResponseWriter, req *http.Request) {

	jdec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	paths := make([]string, 0)

	err := jdec.Decode(&paths)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(fmt.Sprintf("could not parse posted json array input: %v", err)))
		if err != nil {
			slog.Error("failed to write unknown driver name msg to %s: %w", req.RemoteAddr, err)
		}
		return
	}

	results := make(map[string]any)

	for i := range paths {
		driver_name, tag_name, ok := strings.Cut(paths[i], "/")
		if !ok {
			results[paths[i]] = nil
			continue
		}

		driver, ok := drivers[driver_name]
		if !ok {
			results[paths[i]] = nil
			continue
		}

		val, err := driver.Read(tag_name)
		if err != nil {
			results[paths[i]] = nil
			continue
		}
		results[paths[i]] = val

	}

	jenc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	err = jenc.Encode(results)
	if err != nil {
		slog.Error("failed to write json encoded multi-value to %s: %w", req.RemoteAddr, err)
	}

}
