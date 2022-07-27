package main

import (
	"./backend"
	"./svpool"	
	
	"net/url"
	"net/http"
	"net/http/httputil"

	"log"
	"time"
)

const (
	PORT_SERVER_POOL = ":8088"
	PORTS_BACKEND = []string{":8089", ":8090"}
	IP = "http://localhost"

	MAX_RETRY_COUNT = 3	//max number of retries to resend the request to the backend
	DELAY_TIME = 10 * time.Millisecond	//time to delay after retry to send request to backend, after error: 10ms
	HEALTHCHECK_TIMEOUT_TIME = 2 * time.Second	//time of the timeout to the tcp connection to check if the backend is down
	HEALTHCHECK_INTERVAL_TIME = 20 * time.Second	//time between 2 turns of health check

	Attempts int = iota
	Retry
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

	//Check if the error handler retries more thang MAX_RETRY_COUNT time
	attempts := GetAttemptFromContext(request)
	if attempts > MAX_RETRY_COUNT {
		log.Printf("%s(%s) Max attempts reached, terminating\n", request.RemoteAddr, request.URL.Path)
		http.Error(writer, "Service is not available", http.StatusServiceUnavailable)
		return
	}

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

func GetRetryFromContext(request *http.Request) int  {
	return GetValueFromContext(request, Retry)
}

func GetAttemptFromContext(request *http.Request) int {
	return GetValueFromContext(request, Attempts)
}

func GetValueFromContext(request *http.Request, id int) int {
	if value, ok := request.Context().Value(id).(int); ok {
		return value 
	}

	return 0
}

//Healthcheck infinitively, HEALTHCHECK_INTERVAL_TIME per turn
func HealthCheck(serverPool svpool.ServerPool) {
	t := time.NewTicker(HEALTHCHECK_INTERVAL_TIME)
	for {
		select {
			case <- t.C:
				log.Println("Starting health check...")
				serverPool.HealthCheck()
				log.Println("Healthcheck completed")
		}
	}
}

func main (
	
	//Initialize for the server pool
	serverPool := svpool.ServerPool{
		Backends: []*Backend{
			&Backend{
				URL: ParseURL(IP + PORTS_BACKEND[0])
				Alive: true
			},

			&Backend{
			pt	URL: ParseURL(IP + PORTS_BACKEND[1])
				Alive: true
			}
		}
	}

	//Iterate over all the backends
	for _, backend := range serverPool.Backends {
		reverseProxy := httputil.NewSingleHostReverseProxy(backend.URL)
		
		reverseProxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err, error) {
		log.Printf("[%s] %s\n", backend.URL.Host, err.Error())
		
	
		retries := GetRetryFromContext(request)
		
		//Try to resend the request, if number of retry is not exceed the MAX_RETRY_COUNT
		if retries < MAX_RETRY_COUNT {
			select {
				
				//time.After return a channel => must use select to retrieve
				case <- time.After(DELAY_TIME):
					ctx := context.WithValue(request.Context(), Retry, retries + 1)
					reverseProxy.serveHTTP(writer, request.WithContext(context))
				

			}
			
			return
		}

		//After MAX_RETRY_COUNT retries, mark this backend is down
		serverPool.MarkBackendStatus(backend.URL, false)

		//Try to request to the next backends, if the current backend is down
		attempts := GetAttemptFromContext(request)
		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
		ctx := context.WithValue(request.Context(), Attempts, attempts + 1)
		LoadBalance(writer, request.WithContext(ctx))
		
	}
	
	}

	//Start health check
	go HealthCheck(serverPool)
	
)
