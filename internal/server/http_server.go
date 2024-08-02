package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nskforward/dnsx/internal/config"
)

func HTTPListenAndServeTLS(cfg config.Config) {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	indexFile := filepath.Join(dir, "index.html")
	iosSettings := filepath.Join(dir, "dns-profile.mobileconfig")

	http.ListenAndServeTLS(":443", cfg.TLS.Cert, cfg.TLS.Key, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			http.ServeFile(w, r, indexFile)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/ios" {
			f, err := os.Open(iosSettings)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "application/x-apple-aspen-config")
			w.Header().Set("Content-Disposition", "attachment; filename=\"dns-profile.mobileconfig\"")
			w.WriteHeader(200)
			io.Copy(w, f)
			return
		}

		http.NotFound(w, r)
	}))
}
