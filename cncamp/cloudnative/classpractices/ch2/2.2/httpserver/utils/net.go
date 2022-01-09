package utils

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

type (
	RolandClientRealAddr struct {
		Addr string
		Host string
		Port string
	}
	
	// ExtractRolandAddr extract request real addr from a http.Request.
	ExtractRolandAddr func(*http.Request) *RolandClientRealAddr
)

func ExtractRemoteAddrWithDirect() ExtractRolandAddr {
	return func(req *http.Request) *RolandClientRealAddr {
		host, port, _ := net.SplitHostPort(req.RemoteAddr)
		addr, _ := net.ResolveTCPAddr("tcp", req.RemoteAddr)
		if addr != nil {
			host, port = addr.IP.To4().String(), strconv.Itoa(addr.Port)
		}
		
		return &RolandClientRealAddr{
			Addr: req.RemoteAddr,
			Host: host, Port: port,
		}
	}
}

func ExtractAddrFromRequest(req *http.Request) *RolandClientRealAddr {
	directAddr := ExtractRemoteAddrWithDirect()(req)
	
	// extracts IP address using x-real-ip header.
	xRealIP := req.Header.Get(HeaderXRealIP)
	if xRealIP != "" {
		if ip := net.ParseIP(directAddr.Host); ip != nil {
			directAddr.Host = xRealIP
		}
	}
	
	// extracts IP address using x-forwarded-for header.
	xForwardedForList := req.Header[HeaderXForwardedFor]
	if len(xForwardedForList) == 0 {
		return directAddr
	}
	
	ips := append(strings.Split(strings.Join(xForwardedForList, ","), ","), directAddr.Host)
	for i := len(ips) - 1; i >= 0; i-- {
		ip := net.ParseIP(strings.TrimSpace(ips[i]))
		if ip == nil {
			// unable to parse, not trust the records
			return directAddr
		}
	}
	
	// All of the IPs are trusted; return first element because it is furthest from server (best effort strategy).
	directAddr.Host = strings.TrimSpace(ips[0])
	
	return directAddr
}