package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-co-op/gocron" // used to schedule health check function
)

// server struct
type Server struct {
	Host string
	Url  string
	Live bool
}

// server pool url
type ServerPool struct {
	Servers []Server
}

// add host to server pool
func (s *ServerPool) Add(hostUrl *url.URL) {
	s.Servers = append(s.Servers, Server{Host: hostUrl.Host, Url: hostUrl.Host, Live: true})
}

// remove host from server pool
func (s *ServerPool) Remove(hostUrl string) {

	for k, v := range s.Servers {
		if v.Url == hostUrl {
			s.Servers = append(s.Servers[:k], s.Servers[k+1:]...)
		}
	}
}

// check servers health and removes unhealthy server
func healthCheck(serverPool *ServerPool) {

	for _, v := range serverPool.Servers {

		if !getHealth(v.Url) {
			serverPool.Remove(v.Url)
		}
	}

	fmt.Println(serverPool)

}

// return true if server is health false otherwise
func getHealth(url string) bool {

	rep, err := http.Head(url)

	if err != nil {
		log.Println(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from : ", r)
		}
	}()

	return rep.StatusCode == 200

}

var serverPool = ServerPool{[]Server{}}

func main() {

	s := gocron.NewScheduler(time.Local)

	s.Every(5).Seconds().Do(func() {
		go healthCheck(&serverPool)

	})

	s.StartAsync()

	time.Sleep(time.Minute * 2)

}
