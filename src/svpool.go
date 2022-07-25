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

	//Loop entier backends to find out an alive backend
	next := sp.NextIndex();
	size := len(sp.Backends)
	end := size + next;		//start from the 'next' and move a full cycle

	for i := next; i < end; ++i {
		idx := i % size		//the index of the backend in sp.Backends

		//If the tested backend is alive
		// => choose the backend
		if s.Backends[idx].Alive {
			if i != next {
				atomic.StoreUint64(&sp.Current, uint64(idx))
			}
			return sp.Backends[idx]
		}
	}

	//Return nil if all the backends are down
	return nil
}
