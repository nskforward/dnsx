package main

import (
	"log/slog"

	"github.com/nskforward/dnsx/internal/config"
	"github.com/nskforward/dnsx/internal/server"
)

func main() {
	cfg := config.MustLoad()
	go server.HTTPListenAndServeTLS(cfg)
	err := server.TCPListenAndServeTLS(cfg)
	if err != nil {
		slog.Error("failed to start", "error", err)
	}
}
