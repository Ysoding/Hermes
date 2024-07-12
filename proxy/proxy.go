package proxy

import (
	"fmt"
	"io"
	"net/http"
)

type Proxy struct {
	ListenAddr string
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
