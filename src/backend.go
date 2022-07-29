package main

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

//Setter for Alive attribute for the backend with avoid race condition
func (backend *Backend) SetAlive(alive bool) {
	backend.Mux.Lock()
	defer backend.Mux.Unlock()

	backend.Alive = alive
}


//Getter for Alive attribute for the backend with avoid race condition
func (backend *Backend) IsAlive() (alive bool){
	backend.Mux.RLock()
	defer backend.Mux.RUnlock()

	alive = backend.Alive
	return
}
