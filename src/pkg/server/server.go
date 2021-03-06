package server

import (
	"doozer/consensus"
	"doozer/store"
	"log"
	"net"
	"os"
)


// ListenAndServe listens on l, accepts network connections, and
// handles requests according to the doozer protocol.
func ListenAndServe(l net.Listener, canWrite chan bool, st *store.Store, p consensus.Proposer, rwsk, rosk string) {
	var w bool
	for {
		c, err := l.Accept()
		if err != nil {
			if err == os.EINVAL {
				break
			}
			if e, ok := err.(*net.OpError); ok && e.Error == os.EINVAL {
				break
			}
			log.Println(err)
			continue
		}

		// has this server become writable?
		select {
		case w = <-canWrite:
			canWrite = nil
		default:
		}

		go serve(c, st, p, w, rwsk, rosk)
	}
}


func serve(nc net.Conn, st *store.Store, p consensus.Proposer, w bool, rwsk, rosk string) {
	c := &conn{
		c:        nc,
		addr:     nc.RemoteAddr().String(),
		st:       st,
		p:        p,
		canWrite: w,
		rwsk:     rwsk,
		rosk:     rosk,
	}

	c.grant("") // start as if the client supplied a blank password
	c.serve()
	nc.Close()
}
