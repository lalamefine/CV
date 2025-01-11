package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkHandler(b *testing.B) {
	loadFilesToCache()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatal(err)
	}

	benchmarkHandler(b, "serveByFileHandler", serveByFileHandler, req)
	benchmarkHandler(b, "serveByMemCachedFileHandler", serveByMemCachedFileHandler, req)
}

func benchmarkHandler(b *testing.B, name string, handlerFunc http.HandlerFunc, req *http.Request) {
	b.Run(name, func(b *testing.B) {
		var totalDuration time.Duration

		for i := 0; i < b.N; i++ {
			rr := httptest.NewRecorder()
			start := time.Now()
			handlerFunc.ServeHTTP(rr, req)
			duration := time.Since(start)
			totalDuration += duration
		}

		averageDuration := totalDuration / time.Duration(b.N)
		b.Logf("%s average response time: %v", name, averageDuration)
	})
}
