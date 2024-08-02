package server

import (
	"crypto/tls"
	"io"
	"log/slog"
	"net"

	"github.com/nskforward/dnsx/internal/config"
)

func TCPListenAndServeTLS(cfg config.Config) error {
	listener, err := tls.Listen("tcp", cfg.Addr, tlsConfig(cfg.TLS.Cert, cfg.TLS.Key))
	if err != nil {
		return err
	}
	defer listener.Close()

	slog.Info("listening", "addr", cfg.Addr, "ssl", "true")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()
		slog.Info("connected", "addr", conn.RemoteAddr())

		go handleConn(conn, cfg.Upstream)
	}
}

func handleConn(conn net.Conn, backendAddr string) {
	backend, err := net.Dial("tcp", backendAddr)
	if err != nil {
		slog.Error("cannot connect to backend", "error", err)
		return
	}
	closeConns := func() {
		backend.Close()
		conn.Close()
	}
	go func() {
		defer closeConns()
		io.Copy(backend, conn)
	}()
	go func() {
		defer closeConns()
		io.Copy(conn, backend)
	}()
}
