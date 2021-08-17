package healthcheck

import (
	"log"
	"net/http"
	"net/url"
)

// server struct
type Server struct {
	Host string
	URl  string
	Live bool
}

// server pool struct
type ServerPool struct {
	Servers []Server
}

// takes url and adds to server pool
func (serverPool *ServerPool) Register(url url.URL) {
	serverPool.Servers = append(serverPool.Servers, Server{Host: url.Host, URl: url.Host, Live: true})
}

// takes url and deregister server from server pool
func (serverPool *ServerPool) Deregister(url url.URL) {

	for k, v := range serverPool.Servers {
		if v.Host == url.Host {
			serverPool.Servers = append(serverPool.Servers[:k], serverPool.Servers[k+1:]...)

		}
	}
}

// check server health
// takes url and returns true if server is healthy otherwise false
func (serverPool *ServerPool) GetHealth(url url.URL) bool {

	rep, err := http.Head(url.Host)

	if err != nil {
		log.Println(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from : ", r)
		}
	}()

	return rep.StatusCode == 200
}
