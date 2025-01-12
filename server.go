package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var servedDir string

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: server <port> <mode> <servedDir>")
		return
	}

	port := os.Args[1]
	mode := os.Args[2]
	servedDir = os.Args[3] // "./docs"

	switch mode {
	case "file":
		http.HandleFunc("/", serveByFileHandler)
		fmt.Print("Direct file serving mode")
	case "mem":
		loadDirToCache(servedDir)
		http.HandleFunc("/", serveByMemCachedFileHandler)
		fmt.Print("Memory cached file serving mode")
	default:
		fmt.Println("Invalid mode. Please use 'file' or 'mem'")
		return
	}

	fmt.Println(" on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var cache = make(map[string][]byte)

func loadDirToCache(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		path := filepath.Join(path, file.Name())
		if file.IsDir() {
			loadDirToCache(path)
		} else {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.TrimPrefix(path, "docs/")
			fmt.Println("Loaded " + path)
			cache[path] = content
		}
	}
}

func serveByFileHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println("" + r.Method + " /" + name)
	if name == "" {
		name = "index"
	}
	if filepath.Ext(name) == "" {
		name = name + ".html"
	}
	content, err := os.ReadFile(filepath.Join(servedDir, name))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	addMimeTypeHeader(w, name)
	w.Write(content)
}

func serveByMemCachedFileHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println("" + r.Method + " /" + name)
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
