package svpool

import (
	"./backend"
	"sync/atomic"
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

	//TODO: 
}
