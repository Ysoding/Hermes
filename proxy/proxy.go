package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

type Proxy struct {
	ListenAddr string
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
	} else {
		p.handleHttp(w, r)
	}
}

func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to destination: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, fmt.Sprintf("Hijacking error: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	go io.Copy(destConn, clientConn)
	io.Copy(clientConn, destConn)
}

func (p *Proxy) handleHttp(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.String()

	client := &http.Client{}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error creating request: %v", err)
		return
	}

	// copy request header
	req.Header = make(http.Header, len(req.Header))
	for k, v := range req.Header {
		req.Header[k] = append(req.Header[k], v...)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error fetching from target: %v", err)
		return
	}
	defer resp.Body.Close()

	// copy response headers to client
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	// copy response body to client
	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Fprintf(w, "Error copying response: %v", err)
		return
	}
}
