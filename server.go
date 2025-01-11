package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	mode := os.Args[1]

	switch mode {
	case "file":
		http.HandleFunc("/", serveByFileHandler)
		fmt.Print("Direct file serving mode")
	case "mem":
		loadFilesToCache()
		http.HandleFunc("/", serveByMemCachedFileHandler)
		fmt.Print("Memory cached file serving mode")
	default:
		fmt.Println("Invalid mode. Please use 'file' or 'mem'")
		return
	}

	fmt.Println(" on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var cache = make(map[string][]byte)

func loadFilesToCache() {
	files, err := os.ReadDir("./web")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			content, err := os.ReadFile(filepath.Join("./web", file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			name := file.Name()
			cache[name] = content
		}
	}
}

func serveByFileHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index"
	}
	if filepath.Ext(name) == "" {
		name = name + ".html"
	}
	content, err := os.ReadFile(filepath.Join("./web", name))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	addMimeTypeHeader(w, name)
	w.Write(content)
}

func serveByMemCachedFileHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index"
	}
	if filepath.Ext(name) == "" {
		name = name + ".html"
	}
	if content, ok := cache[name]; ok {
		addMimeTypeHeader(w, name)
		w.Header().Set("Cache-Control", "max-age=1600")
		w.Write(content)
	} else {
		http.NotFound(w, r)
	}
}

func addMimeTypeHeader(w http.ResponseWriter, name string) {
	ext := filepath.Ext(name)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}
}
