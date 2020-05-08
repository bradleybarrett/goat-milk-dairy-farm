package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
)

func main() {

	/* Parse command lines args. */
	healthAddrPtr := flag.String("health-addr", "noHealthAddressProvided", "Address of health endpoint (ex. host:port)")
	healthURLPathPtr := flag.String("health-url-path", "/health", "Health URL path of container to register")
	consulAddressPtr := flag.String("consul-addr", "noConsulAddressProvided", "Address of consul instance (ex. host:port)")
	registrationNamePtr := flag.String("registration-name", "noRegistrationNameProvided", "Service name of consul registration")
	registrationAddrPtr := flag.String("registration-addr", "noRegistrationAddressProvided", "Address of consul registration")
	registrationTTLPtr := flag.Int("registration-ttl", 10, "Time-to-live of consul registration (in seconds)")
	flag.Parse()

	/* Register service. */
	register(*healthAddrPtr, *healthURLPathPtr, *consulAddressPtr, *registrationNamePtr, *registrationAddrPtr, time.Duration(*registrationTTLPtr)*time.Second)
}

func register(healthAddr string, healthURLPath string, consulAddress string, registrationName string, registrationAddr string, ttl time.Duration) (*Service, error) {
	s := new(Service)
	s.HealthURL = fmt.Sprintf("%s%s", healthAddr, healthURLPath) // example: 10.0.0.4:8404/health
	s.Name = registrationName
	s.ID = fmt.Sprintf("%s-%s", registrationName, registrationAddr)
	s.TTL = ttl

	/* Check that the lb is healthy. */
	initRetryTime := time.Duration(5) * time.Second
	ok := false
	for !ok {
		ok, _ = s.check()
		if !ok {
			log.Printf("Service %s is DOWN, will check again in %s\n", s.ID, initRetryTime)
			time.Sleep(initRetryTime)
		}
	}
	log.Printf("Service %s is UP, proceed to register the service with consul at %s.\n", s.ID, consulAddress)

	/* Get consul agent with the default config. */
	consulConfigPtr := consul.DefaultConfig()
	(*consulConfigPtr).Address = consulAddress
	c, err := consul.NewClient(consulConfigPtr)
	if err != nil {
		return nil, err
	}
	s.ConsulAgent = c.Agent()

	/* Create the service definition to be registered with consul. */
	registrationHostPortArray := strings.Split(registrationAddr, ":")
	registrationHost := registrationHostPortArray[0]
	registrationPort, err := strconv.Atoi(registrationHostPortArray[1])
	if err != nil {
		registrationPort = 0
		log.Printf("registrationAddr %s is not in host:port format, using default port 0 for consul registration.\n", registrationAddr)
	}
	serviceDef := &consul.AgentServiceRegistration{
		Name:    s.Name,
		ID:      s.ID,
		Address: registrationHost,
		Port:    registrationPort,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}

	/* Register the service definition with consul. */
	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}

	/* Check the service health and repeatedly send heartbeats to consul with the registration status. */
	s.updateTTL(s.check)

	/* This line should not be reached since updateTTL should loop indefinitely. */
	return s, nil
}
