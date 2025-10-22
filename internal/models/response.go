package models

import (
	"time"
)

// Response represents an HTTP response
type Response struct {
	ID         string            `json:"id"`
	RequestID  string            `json:"request_id"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Size       int64             `json:"size"`
	Duration   time.Duration     `json:"duration"`
	CreatedAt  time.Time         `json:"created_at"`
}

// ResponseInfo represents response metadata
type ResponseInfo struct {
	StatusText string            `json:"status_text"`
	Headers    map[string]string `json:"headers"`
	Cookies    []Cookie          `json:"cookies"`
	Timing     ResponseTiming    `json:"timing"`
}

// Cookie represents an HTTP cookie
type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	Secure   bool      `json:"secure"`
	HTTPOnly bool      `json:"http_only"`
}

// ResponseTiming represents timing information
type ResponseTiming struct {
	DNSLookup        time.Duration `json:"dns_lookup"`
	TCPConnection    time.Duration `json:"tcp_connection"`
	TLSHandshake     time.Duration `json:"tls_handshake"`
	ServerProcessing time.Duration `json:"server_processing"`
	ContentTransfer  time.Duration `json:"content_transfer"`
	Total            time.Duration `json:"total"`
}
