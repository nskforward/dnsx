package server

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

type Server struct {
	tcp *dns.Server
	udp *dns.Server
}

func (s *Server) ListenAndServe(addr string, handler func(dns.ResponseWriter, *dns.Msg)) {
	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", handler)
	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", handler)

	s.tcp = &dns.Server{Addr: addr,
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.udp = &dns.Server{Addr: addr,
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := s.tcp.ListenAndServe(); err != nil {
			slog.Error("tcp server start failed", "error", err)
		}
		wg.Done()
	}()
	go func() {
		if err := s.udp.ListenAndServe(); err != nil {
			slog.Error("udp server start failed", "error", err)
		}
		wg.Done()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	slog.Info("server listening", "addr", addr)

	<-sig

	s.tcp.Shutdown()
	s.udp.Shutdown()

	wg.Wait()
}
