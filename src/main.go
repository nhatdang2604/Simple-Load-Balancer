package main

import (
	"./backend"
	"./svpool"	
	"net/http"
	"net/url"
	"net/http/httputil"
)

const (
	PORT_SERVER_POOL = ":8088"
	PORTS_BACKEND = []string{":8089", ":8090"}
	IP = "http://localhost"
)

var (
	serverPool = svpool.ServerPool

	server = http.Server {
		Addr:		PORT
		Handler:	http.HanlderFunc(loadBalance)
	}
)

//Load balances the incoming request
func LoadBalance(writer http.ResponseWriter, request *http.Request) {
	selectedBackend := serverPool.GetNextBackend()
	
	//Throw the error, if all the backends is down
	if nil == selectedBackend {
		http.Error(writer, "Service is not available", http.StatusServiceUnavailable)	
		return
	}

	//Else, pass the request to the backend
	selectedBackend.ReverseProxy.ServeHTTP(writer, request)
}

//Parse the ip string to the url, but with only 1 return value
func ParseURL(u string)(res *url.URL) {
	temp, err := url.Parse(u)
	res = &temp
	if nil != err {
		res = nil
	}

	return
}

func main (
	
	serverPool := ServerPool{
		Backends: []*Backend{
			Backend{
				URL: ParseURL(IP + PORTS_BACKEND[0])
				Alive: true
			}
		}
	}

	//Iterate over all the backends
	for _, backend := range serverPool.Backends {
		reverseProxy := httputil.NewSingleHostReverseProxy(backend.URL)
	}
	
)
