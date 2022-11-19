package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http/httptrace"
	"time"
)

type tracer struct {
	tDNS  time.Time
	dDNS  time.Duration
	tTCP  time.Time
	dTCP  time.Duration
	tTLS  time.Time
	dTLS  time.Duration
	tHTTP time.Time
	dHTTP time.Duration
}

func (t *tracer) GetConn(hostPort string) {
	// fmt.Printf("Get Conn: %+v\n", hostPort)
}

func (t *tracer) DNSStart(dnsInfo httptrace.DNSStartInfo) {
	t.tDNS = time.Now()
}

func (t *tracer) DNSDone(dnsInfo httptrace.DNSDoneInfo) {
	t.dDNS = time.Since(t.tDNS)
	fmt.Printf("dns : %s\n", t.dDNS)
}

func (t *tracer) ConnectStart(network, addr string) {
	t.tTCP = time.Now()
}

func (t *tracer) ConnectDone(network, addr string, err error) {
	t.dTCP = time.Since(t.tTCP)
	fmt.Printf("tcp : %s\n", t.dTCP)
}

func (t *tracer) TLSHandshakeStart() {
	t.tTLS = time.Now()
}

func (t *tracer) TLSHandshakeDone(info tls.ConnectionState, err error) {
	t.dTLS = time.Since(t.tTLS)
	fmt.Printf("tls : %s\n", t.dTCP)
}

func (t *tracer) GotConn(connInfo httptrace.GotConnInfo) {
	t.tHTTP = time.Now()
}

func (t *tracer) GotFirstResponseByte() {
	t.dHTTP = time.Since(t.tHTTP)
	fmt.Printf("htt : %s\n", t.dHTTP)
}

func (t *tracer) Context(ctx context.Context) context.Context {
	trace := &httptrace.ClientTrace{
		GetConn:              t.GetConn,
		DNSStart:             t.DNSStart,
		DNSDone:              t.DNSDone,
		ConnectStart:         t.ConnectStart,
		ConnectDone:          t.ConnectDone,
		TLSHandshakeStart:    t.TLSHandshakeStart,
		TLSHandshakeDone:     t.TLSHandshakeDone,
		GotConn:              t.GotConn,
		GotFirstResponseByte: t.GotFirstResponseByte,
	}
	return httptrace.WithClientTrace(ctx, trace)
}
