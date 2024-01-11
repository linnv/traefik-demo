package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	consul "github.com/hashicorp/consul/api"
)

//Client provides an interface for getting data out of Consul
type Client interface {
	// Get a Service from consul
	Service(string, string) ([]string, error)
	// Register a service with local agent
	Register(string, int) error
	// Deregister a service with local agent
	DeRegister(string) error
}

type client struct {
	consul *consul.Client
}

//NewConsul returns a Client interface for given consul address
func NewConsulClient(addr string) (*client, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	c, err := consul.NewClient(config)
	if err != nil {
		return &client{}, err
	}
	return &client{consul: c}, nil
}

const PrefixRouter = "myservice"

// Register a service with consul local agent - note the tags to define path-prefix is to be used.
func (c *client) Register(id, name, host string, port int, path, health string) error {
	reg := &consul.AgentServiceRegistration{
		ID:   fmt.Sprintf("%s%s%d", id, name, port),
		Name: name,
		// ID: PrefixRouter + "-1",
		// Name:    PrefixRouter,
		Port:    port,
		Address: host,
		Check: &consul.AgentServiceCheck{
			CheckID:       id,
			Name:          "HTTP API health",
			HTTP:          health,
			TLSSkipVerify: true,
			Method:        "GET",
			Interval:      "10s",
			Timeout:       "1s",
		},
		// // Tags: []string{"traefik.enable=true", "traefik.http.routers.myService.rule=Path(`/myservice`)"},
		Tags: []string{
			"traefik.enable=true",
			// 	// Tags:    []string{"traefik.enable=true",},
			// 	// "traefik.http.routers.myService.rule=PathPrefix(`/sdm/myservice`)",
			// 	// "traefik.http.routers.myService.rule=PathPrefix(`/api/greeting`)",

			//same routers name can only set to one rule, no duplicated rule
			// you can set new rule on new router(with new router name)
			// "traefik.http.routers." + name + ".rule=PathPrefix(`/svcwhoami/`)",
			// "traefik.http.routers." + name + "1.rule=PathPrefix(`/svcwhoami/`)",
			// "traefik.http.routers." + name + ".rule=PathPrefix(`/sdm/monkey/`)",
			"traefik.http.routers." + name + "1.rule=PathPrefix(`/sdm/monkey/`)",

			// 	// // "traefik.backend=" + name,
			// 	// "traefik.backend=" + "whoami",
			// 	// // "traefik.frontend.rule=PathPrefix:" + path,
			// // "traefik.http.routers.whoami.rule=Host:whoami.docker.localhost",
			// 	// // "traefik.frontend.rule=PathPrefix(" + path + ")",
			// 	// // "traefik.http.routers.whoami.rule=Host(`whoami.docker.localhost`)",
			// 	//
			// 	// "traefik.frontends.foo.rule=Host:whoami.docker.localhost",
			// 	// "traefik.frontends.bar.rule=PathPrefixStrip:/api/greeting",
			// 	// "traefik.http.routers.myService.rule=Path(`/myservice`)",
		},
	}
	fmt.Printf("reg: %+v\n", reg)
	fmt.Printf("reg: %+v\n", reg.Check)
	return c.consul.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func (c *client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}

// Service return a service
func (c *client) Service(service, tag string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	passingOnly := true
	addrs, meta, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
}

func main() {
	consul := flag.String("consul", "localhost:8500", "Consul host")
	port := flag.Int("port", 10101, "this service port")
	flag.Parse()

	hostname, _ := os.Hostname()
	log.Println("Starting up... ", hostname, " consul host", *consul, " service  ", *port)

	ipAddrs, err := getAllIPAddresses()
	if err != nil {
		panic(err)
	}

	consulClient, _ := NewConsulClient(*consul)
	// Register service with each IP
	for _, hostname := range ipAddrs {
		// registration := newRegistration(ip)
		// err := client.Agent().ServiceRegister(registration)
		// if err != nil {
		//   panic(err)
		// }

		_ = hostname
		// Register Service
		id := fmt.Sprintf("greeting-%v-%v", hostname, *port)
		// // fmt.Printf("id: %+v\n", id)
		health := fmt.Sprintf("http://%v:%v/api/greeting/v1/health", hostname, *port)
		// fmt.Printf("health: %s\n", health)
		_, _ = id, health
		consulClient.Register(id, "svcwhoami", hostname, *port, "/api/greeting", health)

	}

	// Register Service
	// hostname = "192.168.1.237"
	hostname = "192.168.1.237"
	// hostname = fmt.Sprintf("%s:%d", hostname, *port)
	id := fmt.Sprintf("greeting-%v-%v", hostname, *port)
	// consulClient, _ := NewConsulClient(*consul)
	health := fmt.Sprintf("http://%v:%v/api/greeting/v1/health", hostname, *port)
	consulClient.Register(id, "svcwhoami", hostname, *port, "/api/greeting", health)

	router := mux.NewRouter().StrictSlash(true)

	// Define Health Endpoint
	router.Methods("GET").Path("/api/greeting/v1/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		str := fmt.Sprintf("{ 'status':'ok', 'host':'%v:%v' }", hostname, *port)
		fmt.Fprintf(w, str)
		log.Println("/api/greeting/v1/health called")
	})

	// The Hello endpoint for the greeting service
	router.Methods("GET").Path("/api/greeting/v1/hello/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		str := fmt.Sprintf("Hello, %q at %v:%v\n", html.EscapeString(r.URL.Path), hostname, *port)
		rt := rand.Intn(100)
		time.Sleep(time.Duration(rt) * time.Millisecond)
		fmt.Fprintf(w, str)
		log.Println(str)
	})

	router.Methods("GET").Path("/svcwhoami/greeting/v1/hello/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bsReq, _ := httputil.DumpRequest(r, true)
		// if errHttp != nil {
		// 	// ctxlog.Errorf("errHttp %s", errHttp)
		// 	// return
		// }
		log.Printf("bsReq: %s\n", bsReq)
		str := fmt.Sprintf("Hello, %s %q at %v:%v\n %s", r.RemoteAddr, html.EscapeString(r.URL.Path), hostname, *port, string(bsReq))
		rt := rand.Intn(100)
		time.Sleep(time.Duration(rt) * time.Millisecond)
		fmt.Fprintf(w, str)
		log.Println(str)
	})

	router.Methods("GET").Path("/sdm/monkey/greeting/v1/hello/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bsReq, _ := httputil.DumpRequest(r, true)
		// if errHttp != nil {
		// 	// ctxlog.Errorf("errHttp %s", errHttp)
		// 	// return
		// }
		log.Printf("bsReq: %s\n", bsReq)
		str := fmt.Sprintf("Hello, %s %q at %v:%v\n %s", r.RemoteAddr, html.EscapeString(r.URL.Path), hostname, *port, string(bsReq))
		rt := rand.Intn(100)
		time.Sleep(time.Duration(rt) * time.Millisecond)
		fmt.Fprintf(w, str)
		log.Println(str)
	})

	// De-register service at shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("Shutting Down...", sig)
			// consulClient.DeRegister(id)
			consulClient.DeRegister(hostname)
			log.Println("Done...Bye")
			os.Exit(0)
		}
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), router))

}

func getAllIPAddresses() ([]string, error) {
	var ipAddrs []string

	// Get interfaces
	// Loop through interfaces and addresses
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Handle err from interfaces

	// Loop through interfaces
	for _, i := range interfaces {

		// Get IP addresses assigned to interface
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ip := getIPFromAddr(addr)

			// Skip localhost
			if ip.IsLoopback() {
				continue
			}

			ipAddrs = append(ipAddrs, ip.String())
		}
	}

	return ipAddrs, nil
}

func getIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}
