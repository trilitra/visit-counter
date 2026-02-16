package main

import (
	"context"
	"net"
	"net/http"
)

type ctxKey string

const ipKey ctxKey = "ip"

func clientIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}

func GetIP(ctx context.Context) (string, bool) {
	v := ctx.Value(ipKey)
	ip, ok := v.(string)
	return ip, ok
}
