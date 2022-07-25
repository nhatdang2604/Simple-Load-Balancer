package main

import (
	"./backend"
	"./svpool"	
	"net/http"
)

const (
	PORT = ":8088"
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

func main (

)
