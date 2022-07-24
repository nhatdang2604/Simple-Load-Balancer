package svpool

import (
	"./backend"
)

type ServerPool struct {
	Backends	[]*Backend
	Current		uint64
}`
