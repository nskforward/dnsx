package main

import (
	"log/slog"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/nskforward/dnsx/internal/config"
	"github.com/nskforward/dnsx/internal/server"
)

func main() {

	cfg := config.MustLoad()

	var s server.Server
	s.ListenAndServe(cfg.Addr, func(w dns.ResponseWriter, req *dns.Msg) {
		defer w.Close()

		ip, ok := IsKnown(req, cfg)
		if ok {
			w.WriteMsg(CreateReply(req, ip, cfg.TTL))
			return
		}

		slog.Info("proxying")

		resp, err := Lookup(req, cfg.Upstream)
		if err != nil {
			resp = &dns.Msg{}
			resp.SetRcode(req, dns.RcodeServerFailure)
			slog.Info("lookup failed", "error", err)
		}

		w.WriteMsg(resp)
	})
}

func CreateReply(req *dns.Msg, ip net.IP, ttl int) *dns.Msg {
	response := &dns.Msg{}
	response.SetReply(req)

	question := req.Question[0]
	head := dns.RR_Header{
		Name:   question.Name,
		Rrtype: question.Qtype,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}

	var line dns.RR
	if question.Qtype == dns.TypeA {
		line = &dns.A{
			Hdr: head,
			A:   ip,
		}
	} else {
		line = &dns.AAAA{
			Hdr:  head,
			AAAA: ip,
		}
	}
	response.Answer = append(response.Answer, line)
	return response
}

func IsKnown(req *dns.Msg, cfg config.Config) (net.IP, bool) {
	question := req.Question[0]
	if question.Qtype != dns.TypeA && question.Qtype != dns.TypeAAAA {
		return nil, false
	}
	for k, v := range cfg.Routes {
		if k == question.Name {
			return net.ParseIP(v), true
		}
	}
	return nil, false
}

func Lookup(req *dns.Msg, upstream string) (*dns.Msg, error) {
	c := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	// r, _, err := c.Exchange(req, "8.8.8.8:53")
	r, _, err := c.Exchange(req, upstream)

	return r, err
}
