package backend

import (
	"net/url"
	"net/http/httputil"
	"sync"
)

type Backend struct {
	URL		*url.URL
	Alive		bool
	Mux		sync.RWMutex
	ReverseProxy	*httputil.ReverseProxy	
}
