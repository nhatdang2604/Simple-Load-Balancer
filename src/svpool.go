package svpool

import (
	"./backend"
	"sync/atomic"
	"net"
)

type ServerPool struct {
	Backends	[]*Backend
	Current		uint64
}

func (sp *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&sp.Current, uint64(1)) % uint64(len(s.Backends)))
}

//Return next alive backend to take a connection
func (sp *ServerPool) GetNextBackend() *Backend {

	//Loop entier backends to find out an alive backend
	next := sp.NextIndex();
	size := len(sp.Backends)
	end := size + next;		//start from the 'next' and move a full cycle

	for i := next; i < end; ++i {
		idx := i % size		//the index of the backend in sp.Backends

		//If the tested backend is alive
		// => choose the backend
		if s.Backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&sp.Current, uint64(idx))
			}
			return sp.Backends[idx]
		}
	}
	
	//Return nil if all the backends are down
	return nil
}

//Set the status of the backend with the given url with the value alive
func (sp *ServerPool) MarkBackendStatus(backendURL *url.URL, alive bool) {
	
	//Find the backend with the given url
	for _, backend := range sp.Backends {
		if backendURL == backend.URL {
			backend.SetAlive(alive)
			return
		}
	}
}

//Check if the backend is alive by establishing a TCP connection
func IsBackendAlive(u *url.URL) bool {
        connection, err := net.DialTimeout("tcp", u.Host, HEALTHCHECK_TIMEOUT_TIME)

        if nil != err {
                log.Println("SIte unreachable, error: ", err)
                return false
        }

        _ = connection.Close()
        return true

}

//Ping all the backends and update the status
func (sp *ServerPool) HealthCheck() {
	for _, backend := range sp.Backends {
		
		alive := IsBackendAlive(backend.URL)
		backend.SetAlive(alive)
		
		status := "up"
		if !alive {
			status = "down"
		}

		log.Printf("%s [%s]\n", backend.URL, status)
	}
}
